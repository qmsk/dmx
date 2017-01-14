package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
	"github.com/qmsk/go-web"
)

type Options struct {
	LibraryPath string `long:"heads-library" value-name:"PATH"`
}

func (options Options) Heads(config *Config) (*Heads, error) {
	var heads = Heads{
		log:     log.WithField("package", "heads"),
		outputs: make(outputMap),
		heads:   make(headMap),
		groups:  make(groupMap),
		events:  new(Events),
	}

	heads.events.init()

	// preload groups
	for groupID, groupConfig := range config.Groups {
		heads.addGroup(groupID, groupConfig)
	}

	for headID, headConfig := range config.Heads {
		if headConfig.Count > 0 {
			var index uint
			for index = 0; index < headConfig.Count; index++ {
				heads.addHead(headID.index(index), headConfig.index(index), headConfig.headType)
			}
		} else {
			heads.addHead(headID, *headConfig, headConfig.headType)
		}
	}

	// once all heads are patched, init groups
	for _, group := range heads.groups {
		group.init()
	}

	return &heads, nil
}

type headMap map[HeadID]*Head

type Heads struct {
	log     *log.Entry
	outputs outputMap
	heads   headMap
	groups  groupMap
	events  *Events
}

func (heads *Heads) output(universe Universe) *Output {
	output := heads.outputs[universe]
	if output == nil {
		output = &Output{
			log: heads.log.WithField("universe", universe),
			dmx: dmx.MakeUniverse(),
		}

		heads.outputs[universe] = output
	}

	return output
}

// Patch output
func (heads *Heads) Output(config OutputConfig, dmxWriter dmx.Writer) {
	heads.output(config.Universe).init(config, dmxWriter)
}

func (heads *Heads) addGroup(id GroupID, config GroupConfig) *Group {
	group := heads.groups[id]

	if group == nil {
		group = &Group{
			id:     id,
			config: config,
			heads:  make(headMap),
			events: heads.events,
		}

		heads.groups[id] = group
	}

	return group
}

func (heads *Heads) group(id GroupID) *Group {
	if group := heads.groups[id]; group == nil {
		return heads.addGroup(id, GroupConfig{})
	} else {
		return group
	}
}

// Patch head
func (heads *Heads) addHead(id HeadID, config HeadConfig, headType *HeadType) *Head {
	var output = heads.output(config.Universe)
	var head = Head{
		id:       id,
		config:   config,
		headType: headType,
		output:   output,
		events:   heads.events,
	}

	// load head parameters
	head.init()

	// map heads
	heads.heads[id] = &head

	// map groups
	for _, groupID := range config.Groups {
		heads.group(groupID).addHead(&head)
	}

	return &head
}

func (heads *Heads) Each(fn func(head *Head)) {
	for _, head := range heads.heads {
		fn(head)
	}
}

// refresh outputs
func (heads *Heads) Refresh() error {
	var refreshErr error

	for _, output := range heads.outputs {
		if err := output.refresh(); err != nil {
			refreshErr = err
		}
	}

	return refreshErr
}

// Web API
type APIHeads map[HeadID]APIHead

func (heads headMap) makeAPI() APIHeads {
	log.Debug("heads:headMap.makeAPI")

	var apiHeads = make(APIHeads)

	for headID, head := range heads {
		apiHeads[headID] = head.makeAPI()
	}
	return apiHeads
}

type headList headMap

func (heads headList) GetREST() (web.Resource, error) {
	log.Debug("heads:headList.GetREST")

	var apiHeads []APIHead

	for _, head := range heads {
		apiHeads = append(apiHeads, head.makeAPI())
	}

	return apiHeads, nil
}

func (headMap headMap) Index(name string) (web.Resource, error) {
	log.Debugln("heads:headMap.Index", name)

	switch name {
	case "":
		return headList(headMap), nil
	default:
		return headMap[HeadID(name)], nil
	}
}

func (headMap headMap) GetREST() (web.Resource, error) {
	log.Debug("heads:headMap.GetREST")

	return headMap.makeAPI(), nil
}

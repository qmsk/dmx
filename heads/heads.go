package heads

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
	"github.com/SpComb/qmsk-web"
)

type Options struct {
	LibraryPath string `long:"heads-library" value-name:"PATH"`
}

func (options Options) Heads(config *Config) (*Heads, error) {
	var heads = Heads{
		log:     log.WithField("package", "heads"),
		outputs: make(map[Universe]*Output),
		heads:   make(headMap),
	}

	for headID, headConfig := range config.Heads {
		if headType, exists := config.HeadTypes[headConfig.Type]; !exists {
			return nil, fmt.Errorf("Invalid Head.Type=%v", headConfig.Type)
		} else {
			heads.addHead(headID, headConfig, headType)
		}

	}

	return &heads, nil
}

type headMap map[string]*Head

type Heads struct {
	log     *log.Entry
	outputs map[Universe]*Output
	heads   headMap
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
func (heads *Heads) Output(universe Universe, dmxWriter dmx.Writer) {
	heads.output(universe).init(dmxWriter)
}

// Patch head
func (heads *Heads) addHead(id string, config HeadConfig, headType *HeadType) *Head {
	var output = heads.output(config.Universe)
	var head = Head{
		id:       id,
		config:   config,
		headType: headType,
		output:   output,
	}

	head.init()

	heads.heads[id] = &head

	return &head
}

func (heads *Heads) Each(fn func(head *Head)) {
	for _, head := range heads.heads {
		fn(head)
	}
}

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
type APIHeads map[string]APIHead

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
		return headMap[name], nil
	}
}

func (headMap headMap) GetREST() (web.Resource, error) {
	log.Debug("heads:headMap.GetREST")

	return headMap.makeAPI(), nil
}

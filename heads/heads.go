package heads

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
)

type Options struct {
	LibraryPath []string `long:"heads-library" value-name:"PATH"`
}

func (options Options) Heads(config *Config) (*Heads, error) {
	var heads = Heads{
		log:     log.WithField("package", "heads"),
		outputs: make(outputMap),
		heads:   make(headMap),
		groups:  make(groupMap),
		presets: make(presetMap),
		events:  new(Events),
	}

	heads.events.init()

	// preload groups
	for groupID, groupConfig := range config.Groups {
		heads.addGroup(groupID, *groupConfig)
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

	//  load presets
	for presetID, presetConfig := range config.Presets {
		if err := heads.addPreset(presetID, *presetConfig); err != nil {
			return nil, fmt.Errorf("Load preset=%v: %v", presetID, err)
		}
	}

	return &heads, nil
}

type Heads struct {
	log     *log.Entry
	outputs outputMap
	heads   headMap
	groups  groupMap
	presets presetMap
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

func (heads *Heads) addPreset(id PresetID, config PresetConfig) error {
	var preset = Preset{
		ID:     id,
		Config: config,
		Groups: make(map[GroupID]APIGroupParams),
	}

	if preset.Config.All != nil {
		preset.initAll(heads.heads)
	}

	for groupID, presetParameters := range preset.Config.Groups {
		if group := heads.groups[groupID]; group == nil {
			return fmt.Errorf("No such group: %v", groupID)
		} else if err := presetParameters.initGroup(group); err != nil {
			return err
		} else {
			preset.initGroup(group, presetParameters)
		}
	}

	for headID, presetParameters := range preset.Config.Heads {
		if head := heads.heads[headID]; head == nil {
			return fmt.Errorf("No such head: %v", headID)
		} else if err := presetParameters.initHead(head); err != nil {
			return err
		} else {
			preset.initHead(head, presetParameters)
		}
	}

	heads.presets[id] = &preset

	return nil
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

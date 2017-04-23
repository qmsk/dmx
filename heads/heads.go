package heads

import (
	"fmt"
	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/logging"
)

type Options struct {
	Log logging.Option `long:"log.heads"`

	LibraryPath []string `long:"heads-library" value-name:"PATH"`
}

func (options Options) Heads(config *Config) (*Heads, error) {
	options.Log.Package = "heads"

	var heads = Heads{
		options: options,
		log:     options.Log.Logger("package", "heads"),

		outputs: make(outputMap),
		heads:   make(headMap),
		groups:  make(groupMap),
		presets: make(presetMap),
		events: &events{
			log: options.Log.Logger("events", nil),
		},
	}

	if err := heads.load(config); err != nil {
		return nil, err
	}

	return &heads, nil
}

type Heads struct {
	options Options
	log     logging.Logger
	config  *Config
	outputs outputMap
	heads   headMap
	groups  groupMap
	presets presetMap
	events  *events // optional
}

func (heads *Heads) output(universe Universe) *Output {
	output := heads.outputs[universe]
	if output == nil {
		output = &Output{
			events: heads.events,
		}

		output.init(heads.log, universe)

		heads.outputs[universe] = output
	}

	return output
}

// Patch output
// XXX: not goroutine-safe...
func (heads *Heads) Output(universe Universe, config OutputConfig, writer dmx.Writer) {
	heads.output(universe).connect(config, writer)
}

func (heads *Heads) load(config *Config) error {
	heads.config = config

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
			return fmt.Errorf("Load preset=%v: %v", presetID, err)
		}
	}

	return nil
}

func (heads *Heads) addGroup(id GroupID, config GroupConfig) *Group {
	group := heads.groups[id]

	if group == nil {
		group = &Group{
			log:    heads.options.Log.Logger("group", id),
			id:     id,
			config: config,
			heads:  make(headMap),
			colors: make(ColorMap),
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
		log:      heads.options.Log.Logger("head", id),
		id:       id,
		config:   config,
		headType: headType,
		output:   output,
		events:   heads.events,
		groups:   make(groupMap),
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
		log:    heads.options.Log.Logger("preset", id),
		events: heads.events,
		ID:     id,
		Config: config,
		Groups: make(map[GroupID]PresetParameters),
		Heads:  make(map[HeadID]PresetParameters),
	}

	if preset.Config.All != nil {
		preset.initAll(heads.heads, heads.groups)
	}

	for groupID, presetParameters := range preset.Config.Groups {
		if group := heads.groups[GroupID(groupID)]; group == nil {
			return fmt.Errorf("No such group: %v", groupID)
		} else {
			preset.initGroup(group, presetParameters)
		}
	}

	for headID, presetParameters := range preset.Config.Heads {
		if head := heads.heads[HeadID(headID)]; head == nil {
			return fmt.Errorf("No such head: %v", headID)
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
		if err := output.Refresh(); err != nil {
			refreshErr = err
		}
	}

	return refreshErr
}

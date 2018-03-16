package heads

import (
	"fmt"
	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
)

func MakeController() Controller {
	return Controller{
		outputs: make(outputMap),
		heads:   make(headMap),
		groups:  make(groupMap),
		presets: make(presetMap),
		events:  &events{},
	}
}

type Controller struct {
	log     logging.Logger
	outputs outputMap
	heads   headMap
	groups  groupMap
	presets presetMap
	events  *events // optional
}

func (controller *Controller) output(universe Universe) *Output {
	output := controller.outputs[universe]
	if output == nil {
		output = &Output{
			events: controller.events,
		}

		output.init(controller.log, universe)

		controller.outputs[universe] = output
	}

	return output
}

// Patch output
// XXX: not goroutine-safe...
func (controller *Controller) Output(universe Universe, config OutputConfig, writer dmx.Writer) {
	controller.output(universe).connect(config, writer)
}

func (controller *Controller) Load(config api.Config) error {
	// preload groups
	for groupID, groupConfig := range config.Groups {
		controller.addGroup(groupID, groupConfig)
	}

	if err := loadHeads(config, func(headID api.HeadID, headConfig api.HeadConfig, headType api.HeadType) error {
		controller.addHead(headID, headConfig, headType)

		return nil
	}); err != nil {
		return err
	}

	// once all heads are patched, init groups
	for _, group := range controller.groups {
		group.init()
	}

	//  load presets
	for presetID, presetConfig := range config.Presets {
		if err := controller.addPreset(presetID, presetConfig); err != nil {
			return fmt.Errorf("Load preset %v: %v", presetID, err)
		}
	}

	return nil
}

func (controller *Controller) addGroup(id GroupID, config GroupConfig) *Group {
	group := controller.groups[id]

	if group == nil {
		group = &Group{
			log:    controller.options.Log.Logger("group", id),
			id:     id,
			config: config,
			heads:  make(headMap),
			colors: make(ColorMap),
			events: controller.events,
		}

		controller.groups[id] = group
	}

	return group
}

func (controller *Controller) group(id GroupID) *Group {
	if group := controller.groups[id]; group == nil {
		return controller.addGroup(id, GroupConfig{})
	} else {
		return group
	}
}

// Patch head
func (controller *Controller) addHead(id HeadID, config HeadConfig, headType *HeadType) *Head {
	var output = controller.output(config.Universe)
	var head = Head{
		log:      controller.options.Log.Logger("head", id),
		id:       id,
		config:   config,
		headType: headType,
		output:   output,
		events:   controller.events,
		groups:   make(groupMap),
	}

	// load head parameters
	head.init()

	// map controller
	controller.heads[id] = &head

	// map groups
	for _, groupID := range config.Groups {
		controller.group(groupID).addHead(&head)
	}

	return &head
}

func (controller *Controller) addPreset(id PresetID, config PresetConfig) error {
	var preset = Preset{
		log:    controller.options.Log.Logger("preset", id),
		events: controller.events,
		ID:     id,
		Config: config,
		Groups: make(map[GroupID]PresetParameters),
		Heads:  make(map[HeadID]PresetParameters),
	}

	if preset.Config.All != nil {
		preset.initAll(controller.heads, controller.groups)
	}

	for groupID, presetParameters := range preset.Config.Groups {
		if group := controller.groups[GroupID(groupID)]; group == nil {
			return fmt.Errorf("No such group: %v", groupID)
		} else {
			preset.initGroup(group, presetParameters)
		}
	}

	for headID, presetParameters := range preset.Config.Heads {
		if head := controller.heads[HeadID(headID)]; head == nil {
			return fmt.Errorf("No such head: %v", headID)
		} else {
			preset.initHead(head, presetParameters)
		}
	}

	controller.presets[id] = &preset

	return nil
}

func (controller *Controller) Each(fn func(head *Head)) {
	for _, head := range controller.heads {
		fn(head)
	}
}

// refresh outputs
func (controller *Controller) Refresh() error {
	var refreshErr error

	for _, output := range controller.outputs {
		if err := output.Refresh(); err != nil {
			refreshErr = err
		}
	}

	return refreshErr
}

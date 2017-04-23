package heads

import (
	"net/http"

	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"

	"github.com/BurntSushi/toml"
)

// Config
type PresetID string

type PresetParameters struct {
	head  *Head
	group *Group

	Intensity *APIIntensity
	Color     *APIColor
}

func (presetParameters PresetParameters) IsZero() bool {
	if presetParameters.Intensity != nil && !presetParameters.Intensity.IsZero() {
		return false
	}
	if presetParameters.Color != nil && !presetParameters.Color.IsZero() {
		return false
	}
	return true
}

// Do the parameters set in this preset override any of those parameters set in the other preset?
func (presetParameters PresetParameters) Overrides(other PresetParameters) bool {
	if presetParameters.Intensity == nil || other.Intensity == nil {

	} else if !presetParameters.Intensity.Equals(*other.Intensity) {
		return true
	}

	if presetParameters.Color == nil || other.Color == nil {

	} else if !presetParameters.Color.Equals(*other.Color) {
		return true
	}

	return false
}

func (presetParameters *PresetParameters) init() error {
	if presetParameters.Intensity != nil {
		if presetParameters.group != nil {
			if err := presetParameters.Intensity.initGroup(presetParameters.group.intensity); err != nil {
				return err
			}
		}
		if presetParameters.head != nil {
			if err := presetParameters.Intensity.initHead(presetParameters.head.parameters.Intensity); err != nil {
				return err
			}
		}
	}

	if presetParameters.Color != nil {
		if presetParameters.group != nil {
			if err := presetParameters.Color.initGroup(presetParameters.group.color); err != nil {
				return err
			}
		}
		if presetParameters.head != nil {
			if err := presetParameters.Color.initHead(presetParameters.head.parameters.Color); err != nil {
				return err
			}
		}
	}

	return nil
}

func (presetParameters PresetParameters) scaleIntensity(scaleIntensity Intensity) PresetParameters {
	if presetParameters.Intensity != nil {
		apiIntensity := *presetParameters.Intensity
		apiIntensity.ScaleIntensity = &scaleIntensity

		presetParameters.Intensity = &apiIntensity
	}

	if presetParameters.Color != nil {
		apiColor := *presetParameters.Color
		apiColor.ScaleIntensity = &scaleIntensity

		presetParameters.Color = &apiColor
	}

	return presetParameters
}

// requires initHead/initGroup
func (presetParameters PresetParameters) Apply() error {
	if presetParameters.Intensity == nil {

	} else if err := presetParameters.Intensity.Apply(); err != nil {
		return err
	}

	if presetParameters.Color == nil {

	} else if err := presetParameters.Color.Apply(); err != nil {
		return err
	}

	return nil
}

type PresetConfig struct {
	Name   string
	All    *PresetParameters
	Groups map[string]PresetParameters
	Heads  map[string]PresetParameters
}

// Heads.Presets
type presetMap map[PresetID]*Preset

func (presetMap presetMap) GetREST() (web.Resource, error) {
	return presetMap, nil
}

func (presetMap presetMap) Index(name string) (web.Resource, error) {
	return presetMap[PresetID(name)], nil
}

type Preset struct {
	log    logging.Logger
	events Events

	ID     PresetID
	Config PresetConfig

	allHeads  headMap
	allGroups groupMap
	Groups    map[GroupID]PresetParameters
	Heads     map[HeadID]PresetParameters
}

func (preset *Preset) initAll(heads headMap, groups groupMap) {
	preset.allHeads = heads
	preset.allGroups = groups
}

func (preset *Preset) initGroup(group *Group, presetParameters PresetParameters) error {
	var groupParameters = PresetParameters{
		group:     group,
		Intensity: presetParameters.Intensity,
		Color:     presetParameters.Color,
	}

	if err := groupParameters.init(); err != nil {
		return err
	}

	preset.Groups[group.id] = groupParameters

	return nil
}

func (preset *Preset) initHead(head *Head, presetParameters PresetParameters) error {
	var headParameters = PresetParameters{
		head:      head,
		Intensity: presetParameters.Intensity,
		Color:     presetParameters.Color,
	}

	if err := headParameters.init(); err != nil {
		return err
	}

	preset.Heads[head.id] = headParameters

	return nil
}

func (preset *Preset) GetREST() (web.Resource, error) {
	return preset, nil
}

func (preset *Preset) PostREST() (web.Resource, error) {
	return &APIPresetParams{preset: preset}, nil
}

// API POST
type APIPresetParams struct {
	preset *Preset

	Intensity *Intensity
}

func (apiPresetParams APIPresetParams) Apply() error {
	var preset = apiPresetParams.preset
	var event APIEvents

	preset.log.Info("Apply")

	if allParams := preset.Config.All; allParams != nil {
		for _, head := range preset.allHeads {
			var headParams = PresetParameters{
				head: head,
			}

			// all params are optional
			if allParams.Intensity != nil && head.parameters.Intensity != nil {
				headParams.Intensity = allParams.Intensity
			}

			if allParams.Color != nil && head.parameters.Color != nil {
				headParams.Color = allParams.Color
			}

			if apiPresetParams.Intensity != nil {
				headParams = headParams.scaleIntensity(*apiPresetParams.Intensity)
			}

			if err := headParams.init(); err != nil {
				return err
			} else if err := headParams.Apply(); err != nil {
				return err
			}
		}

		// update everything
		event.addHeads(preset.allHeads)
		event.addGroups(preset.allGroups)

		// also update groups after heads have been updated
		for _, group := range apiPresetParams.preset.allGroups {
			if err := group.Apply(); err != nil {
				return err
			}
		}
	}

	for _, apiGroupParams := range apiPresetParams.preset.Groups {
		if apiPresetParams.Intensity != nil {
			apiGroupParams = apiGroupParams.scaleIntensity(*apiPresetParams.Intensity)
		}

		if err := apiGroupParams.Apply(); err != nil {
			return err
		}

		event.addGroup(apiGroupParams.group)
		event.addHeads(apiGroupParams.group.heads)
	}

	for _, apiHeadParams := range apiPresetParams.preset.Heads {
		if apiPresetParams.Intensity != nil {
			apiHeadParams = apiHeadParams.scaleIntensity(*apiPresetParams.Intensity)
		}

		if err := apiHeadParams.Apply(); err != nil {
			return err
		}

		event.addHead(apiHeadParams.head)
		event.addGroups(apiHeadParams.head.groups)
	}

	preset.events.update(event)

	return nil
}

// GET /config/preset.toml
type httpConfigPreset struct {
	heads *Heads
}

// Export a preset configuration from the current state
func (heads *Heads) ConfigPreset() PresetConfig {
	var allParameters = PresetParameters{
		Intensity: &APIIntensity{},
		Color:     &APIColor{},
	}

	var presetConfig = PresetConfig{
		All:    &allParameters,
		Groups: make(map[string]PresetParameters),
		Heads:  make(map[string]PresetParameters),
	}

	for groupID, group := range heads.groups {
		var presetParameters = PresetParameters{
			Intensity: group.intensity.makeAPI(),
			Color:     group.color.makeAPI(),
		}

		if presetParameters.Overrides(allParameters) {

		} else {
			continue
		}

		presetConfig.Groups[string(groupID)] = presetParameters
	}

	for headID, head := range heads.heads {
		var presetParameters = PresetParameters{
			Intensity: head.parameters.Intensity.makeAPI(),
			Color:     head.parameters.Color.makeAPI(),
		}

		var baseParameters = allParameters

		for groupID, _ := range head.groups {
			if groupParameters, exists := presetConfig.Groups[string(groupID)]; exists {
				baseParameters = groupParameters
			}
		}

		if presetParameters.Overrides(baseParameters) {

		} else {
			continue
		}

		presetConfig.Heads[string(headID)] = presetParameters
	}

	return presetConfig
}

func (httpConfigPreset httpConfigPreset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/toml")

	var presetConfig = httpConfigPreset.heads.ConfigPreset()

	if err := toml.NewEncoder(w).Encode(presetConfig); err != nil {
		panic(err)
	}
}

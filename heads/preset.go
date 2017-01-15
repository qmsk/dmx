package heads

import (
	"github.com/qmsk/go-web"
)

// Config
type PresetID string

type PresetParameters struct {
	resource  web.MutableResource
	Intensity *APIIntensity
	Color     *APIColor
}

// XXX: this is called on a copy for the error-checking
func (presetParameters *PresetParameters) initGroup(group *Group) error {
	if presetParameters.Intensity != nil {
		if err := presetParameters.Intensity.initGroup(group.intensity); err != nil {
			return err
		}
	}

	if presetParameters.Color != nil {
		if err := presetParameters.Color.initGroup(group.color); err != nil {
			return err
		}
	}

	return nil
}

// XXX: this is called on a copy for the error-checking
func (presetParameters *PresetParameters) initHead(head *Head) error {
	if presetParameters.Intensity != nil {
		if err := presetParameters.Intensity.initHead(head.parameters.Intensity); err != nil {
			return err
		}
	}

	if presetParameters.Color != nil {
		if err := presetParameters.Color.initHead(head.parameters.Color); err != nil {
			return err
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
	Groups map[GroupID]PresetParameters
	Heads  map[HeadID]PresetParameters
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
		resource:  group,
		Intensity: presetParameters.Intensity,
		Color:     presetParameters.Color,
	}

	if err := groupParameters.initGroup(group); err != nil {
		return err
	}

	preset.Groups[group.id] = groupParameters

	return nil
}

func (preset *Preset) initHead(head *Head, presetParameters PresetParameters) error {
	var headParameters = PresetParameters{
		resource:  head,
		Intensity: presetParameters.Intensity,
		Color:     presetParameters.Color,
	}

	if err := headParameters.initHead(head); err != nil {
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
	if allParams := apiPresetParams.preset.Config.All; allParams != nil {
		for _, head := range apiPresetParams.preset.allHeads {
			var headParams = PresetParameters{}

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

			if err := headParams.initHead(head); err != nil {
				return err
			} else if err := headParams.Apply(); err != nil {
				return err
			} else if err := head.Apply(); err != nil {
				return err
			}
		}

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
		} else if err := apiGroupParams.resource.Apply(); err != nil {
			return err
		}
	}

	for _, apiHeadParams := range apiPresetParams.preset.Heads {
		if apiPresetParams.Intensity != nil {
			apiHeadParams = apiHeadParams.scaleIntensity(*apiPresetParams.Intensity)
		}

		if err := apiHeadParams.Apply(); err != nil {
			return err
		} else if err := apiHeadParams.resource.Apply(); err != nil {
			return err
		}
	}

	return nil
}

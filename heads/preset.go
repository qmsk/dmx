package heads

import (
	"github.com/qmsk/go-web"
)

// Config
type PresetID string

type PresetParameters struct {
	Intensity *APIIntensity
	Color     *APIColor
}

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

type PresetConfig struct {
	Name   string
	All    *PresetParameters
	Groups map[GroupID]*PresetParameters
	Heads  map[HeadID]*PresetParameters
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
}

func (preset *Preset) GetREST() (web.Resource, error) {
	return preset, nil
}

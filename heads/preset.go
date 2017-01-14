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

	all    headMap
	Groups map[GroupID]APIGroupParams
	Heads  map[HeadID]APIHeadParams
}

func (preset *Preset) initAll(heads headMap) {
	preset.all = heads
}

func (preset *Preset) initGroup(group *Group, presetParameters PresetParameters) {
	preset.Groups[group.id] = APIGroupParams{
		group:     group,
		Intensity: presetParameters.Intensity,
		Color:     presetParameters.Color,
	}
}

func (preset *Preset) initHead(head *Head, presetParameters PresetParameters) {
	preset.Heads[head.id] = APIHeadParams{
		head:      head,
		Intensity: presetParameters.Intensity,
		Color:     presetParameters.Color,
	}
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
}

func (apiPresetParams APIPresetParams) Apply() error {
	for _, head := range apiPresetParams.preset.all {
		var headParams = APIHeadParams{
			head:      head,
			Intensity: apiPresetParams.preset.Config.All.Intensity,
			Color:     apiPresetParams.preset.Config.All.Color,
		}

		if err := headParams.Apply(); err != nil {
			return err
		} else if err := head.Apply(); err != nil {
			return err
		}
	}

	for _, apiGroupParams := range apiPresetParams.preset.Groups {
		if err := apiGroupParams.Apply(); err != nil {
			return err
		} else if err := apiGroupParams.group.Apply(); err != nil {
			return err
		}
	}

	for _, apiHeadParams := range apiPresetParams.preset.Heads {
		if err := apiHeadParams.Apply(); err != nil {
			return err
		} else if err := apiHeadParams.head.Apply(); err != nil {
			return err
		}
	}

	return nil
}

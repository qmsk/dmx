package api

type PresetParameterConfig struct {
	Intensity *Intensity
	Color     *Color
}

type PresetConfig struct {
	Name   string
	All    *PresetParameterConfig
	Groups map[string]PresetParameterConfig
	Heads  map[string]PresetParameterConfig
}

type PresetID string
type Presets map[PresetID]Preset

type Preset struct {
	ID     PresetID
	Config PresetConfig

	Groups map[GroupID]PresetParameterConfig
	Heads  map[HeadID]PresetParameterConfig
}

type PresetParams struct {
	Intensity *Intensity
}

package api

type PresetConfigParameters struct {
	Intensity *Intensity
	Color     *Color
}

type PresetGroups map[GroupID]PresetConfigParameters
type PresetHeads map[HeadID]PresetConfigParameters

type PresetConfig struct {
	Name   string
	All    *PresetConfigParameters
	Groups PresetGroups
	Heads  PresetHeads
}

type PresetID string
type Presets map[PresetID]Preset

type Preset struct {
	ID     PresetID
	Config PresetConfig

	Groups PresetGroups
	Heads  PresetHeads
}

type PresetParams struct {
	Intensity *Value // TODO: rename to Scale
}

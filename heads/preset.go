package heads

type PresetID string

type PresetParameters struct {
	Intensity *APIIntensity
	Color     *APIColor
}

type PresetConfig struct {
	All    *PresetParameters
	Groups map[GroupID]PresetParameters
	Heads  map[GroupID]PresetParameters
}

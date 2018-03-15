package api

type Intensity struct {
	Intensity Value
}

type IntensityParams struct {
	ScaleIntensity *Value
	Intensity
}

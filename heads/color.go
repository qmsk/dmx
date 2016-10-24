package heads

import "github.com/SpComb/qmsk-web"

type ColorRGB struct {
	R, G, B Value
}

// Linear RGB intensity scaling
func (color ColorRGB) scaleIntensity(intensity Intensity) ColorRGB {
	return ColorRGB{
		R: color.R * Value(intensity),
		G: color.G * Value(intensity),
		B: color.B * Value(intensity),
	}
}

type HeadColor struct {
	red       *Channel
	green     *Channel
	blue      *Channel
	intensity *Channel
}

func (it HeadColor) Exists() bool {
	return it.red != nil || it.green != nil || it.blue != nil
}

func (hc HeadColor) SetRGB(colorRGB ColorRGB) {
	if hc.red != nil {
		hc.red.SetValue(colorRGB.R)
	}
	if hc.green != nil {
		hc.green.SetValue(colorRGB.G)
	}
	if hc.blue != nil {
		hc.blue.SetValue(colorRGB.B)
	}
}

// Set color with intensity, using either head intensity channel or linear RGB scaling
func (hc HeadColor) SetRGBIntensity(colorRGB ColorRGB, intensity Intensity) {
	if hc.intensity != nil {
		hc.SetRGB(colorRGB)
		hc.intensity.SetValue(Value(intensity))
	} else {
		hc.SetRGB(colorRGB.scaleIntensity(intensity))
	}
}

// Web API
type APIHeadColor struct {
	Red   Value
	Green Value
	Blue  Value
}

func (headColor HeadColor) makeAPI() *APIHeadColor {
	if !headColor.Exists() {
		return nil
	}

	return &APIHeadColor{}
}

func (headColor HeadColor) GetREST() (web.Resource, error) {
	return headColor.makeAPI(), nil
}

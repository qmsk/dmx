package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-web"
)

type ColorRGB struct {
	Red   Value
	Green Value
	Blue  Value
}

// Linear RGB intensity scaling
func (color ColorRGB) scaleIntensity(intensity Intensity) ColorRGB {
	return ColorRGB{
		Red:   color.Red * Value(intensity),
		Green: color.Green * Value(intensity),
		Blue:  color.Blue * Value(intensity),
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

func (hc HeadColor) GetRGB() (colorRGB ColorRGB) {
	if hc.red != nil {
		colorRGB.Red = hc.red.GetValue()
	}
	if hc.green != nil {
		colorRGB.Green = hc.green.GetValue()
	}
	if hc.blue != nil {
		colorRGB.Blue = hc.blue.GetValue()
	}
	return
}

func (hc HeadColor) SetRGB(colorRGB ColorRGB) ColorRGB {
	if hc.red != nil {
		colorRGB.Red = hc.red.SetValue(colorRGB.Red)
	}
	if hc.green != nil {
		colorRGB.Green = hc.green.SetValue(colorRGB.Green)
	}
	if hc.blue != nil {
		colorRGB.Blue = hc.blue.SetValue(colorRGB.Blue)
	}
	return colorRGB
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
	headColor HeadColor

	ColorRGB
}

func (headColor HeadColor) makeAPI() *APIHeadColor {
	if !headColor.Exists() {
		return nil
	}

	return &APIHeadColor{
		headColor: headColor,
		ColorRGB:  headColor.GetRGB(),
	}
}

func (headColor HeadColor) GetREST() (web.Resource, error) {
	return headColor.makeAPI(), nil
}
func (headColor HeadColor) PostREST() (web.Resource, error) {
	return headColor.makeAPI(), nil
}

func (apiHeadColor *APIHeadColor) Apply() error {
	log.Debugln("heads:APIHeadColor.Apply", apiHeadColor.ColorRGB)

	apiHeadColor.ColorRGB = apiHeadColor.headColor.SetRGB(apiHeadColor.ColorRGB)

	return nil
}

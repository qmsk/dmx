package heads

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/qmsk/go-web"
)

// Config Channels
type ColorChannel string

const (
	ColorChannelRed   = "red"
	ColorChannelGreen = "green"
	ColorChannelBlue  = "blue"
)

// Config
type ColorID string
type ColorMap map[ColorID]ColorRGB

func (colorMap ColorMap) Merge(mergeMap ColorMap) {
	for colorID, color := range mergeMap {
		colorMap[colorID] = color
	}
}

// Types
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

// Head.Color
type HeadColor struct {
	red       *Channel
	green     *Channel
	blue      *Channel
	intensity *Channel
}

func (it HeadColor) exists() bool {
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

func (headColor *HeadColor) makeAPI() *APIColor {
	if headColor == nil {
		return nil
	}

	return &APIColor{
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

// Group.Color
type GroupColor struct {
	headColors map[HeadID]HeadColor
}

func (groupColor GroupColor) exists() bool {
	return len(groupColor.headColors) > 0
}

// Return one color for the group
func (groupColor GroupColor) GetRGB() (colorRGB ColorRGB) {
	for _, headColor := range groupColor.headColors {
		// This works fine assuming they are all the same color :)
		return headColor.GetRGB()
	}

	return
}

func (groupColor GroupColor) SetRGB(colorRGB ColorRGB) ColorRGB {
	for _, headColor := range groupColor.headColors {
		headColor.SetRGB(colorRGB)
	}

	return colorRGB
}

func (groupColor *GroupColor) makeAPI() *APIColor {
	if groupColor == nil {
		return nil
	}

	return &APIColor{
		groupColor: groupColor,
		ColorRGB:   groupColor.GetRGB(),
	}
}

// Web API
type APIColor struct {
	headColor  *HeadColor
	groupColor *GroupColor

	ColorRGB
}

func (apiColor *APIColor) initHead(headColor *HeadColor) error {
	if headColor == nil {
		return fmt.Errorf("Head does not support color")
	}

	apiColor.headColor = headColor

	return nil
}

func (apiColor *APIColor) initGroup(groupColor *GroupColor) error {
	if groupColor == nil {
		return fmt.Errorf("Group does not support color")
	}

	apiColor.groupColor = groupColor

	return nil
}

func (apiColor *APIColor) Apply() error {
	if apiColor.headColor != nil {
		log.Debugln("heads:APIColor.Apply head", apiColor.headColor, apiColor.ColorRGB)

		apiColor.ColorRGB = apiColor.headColor.SetRGB(apiColor.ColorRGB)
	}

	if apiColor.groupColor != nil {
		log.Debugln("heads:APIColor.Apply group", apiColor.groupColor, apiColor.ColorRGB)

		apiColor.ColorRGB = apiColor.groupColor.SetRGB(apiColor.ColorRGB)
	}

	return nil
}

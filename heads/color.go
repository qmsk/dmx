package heads

import (
	"fmt"
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

type ColorMap map[ColorID]Color

// Merge in new colors from given map
// Preserves existing colors
func (colorMap ColorMap) Merge(mergeMap ColorMap) {
	for colorID, color := range mergeMap {
		if _, exists := colorMap[colorID]; !exists {
			colorMap[colorID] = color
		}
	}
}

// Types
type Color struct {
	Red   Value
	Green Value
	Blue  Value
}

func (color Color) IsZero() bool {
	return color.Red == 0.0 && color.Green == 0.0 && color.Blue == 0.0
}

// Linear RGB intensity scaling
func (color Color) ScaleIntensity(intensity Intensity) Color {
	return Color{
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

func (hc HeadColor) Get() (color Color) {
	if hc.red != nil {
		color.Red = hc.red.GetValue()
	}
	if hc.green != nil {
		color.Green = hc.green.GetValue()
	}
	if hc.blue != nil {
		color.Blue = hc.blue.GetValue()
	}
	return
}

func (hc HeadColor) Set(color Color) Color {
	if hc.red != nil {
		color.Red = hc.red.SetValue(color.Red)
	}
	if hc.green != nil {
		color.Green = hc.green.SetValue(color.Green)
	}
	if hc.blue != nil {
		color.Blue = hc.blue.SetValue(color.Blue)
	}
	return color
}

// Set color with intensity, using either head intensity channel or linear RGB scaling
func (hc HeadColor) SetIntensity(color Color, intensity Intensity) {
	if hc.intensity != nil {
		hc.Set(color)
		hc.intensity.SetValue(Value(intensity))
	} else {
		hc.Set(color.ScaleIntensity(intensity))
	}
}

func (headColor *HeadColor) makeAPI() *APIColor {
	if headColor == nil {
		return nil
	}

	return &APIColor{
		headColor: headColor,
		Color:     headColor.Get(),
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
func (groupColor GroupColor) Get() (color Color) {
	for _, headColor := range groupColor.headColors {
		// This works fine assuming they are all the same color :)
		return headColor.Get()
	}

	return
}

func (groupColor GroupColor) Set(color Color) Color {
	for _, headColor := range groupColor.headColors {
		headColor.Set(color)
	}

	return color
}

func (groupColor *GroupColor) makeAPI() *APIColor {
	if groupColor == nil {
		return nil
	}

	return &APIColor{
		groupColor: groupColor,
		Color:      groupColor.Get(),
	}
}

// Web API
type APIColor struct {
	headColor  *HeadColor
	groupColor *GroupColor

	ScaleIntensity *Intensity
	Color
}

func (apiColor APIColor) IsZero() bool {
	return apiColor.Color.IsZero()
}
func (apiColor APIColor) Equals(other APIColor) bool {
	return apiColor.Color == other.Color
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
	if apiColor.ScaleIntensity != nil {
		apiColor.Color = apiColor.Color.ScaleIntensity(*apiColor.ScaleIntensity)
	}

	if apiColor.headColor != nil {
		apiColor.Color = apiColor.headColor.Set(apiColor.Color)
	}

	if apiColor.groupColor != nil {
		apiColor.Color = apiColor.groupColor.Set(apiColor.Color)
	}

	return nil
}

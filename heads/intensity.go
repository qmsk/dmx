package heads

import (
	"fmt"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/go-web"
)

// Head.Intensity
type HeadIntensity struct {
	channel *Channel
}

func (it HeadIntensity) exists() bool {
	return it.channel != nil
}

func (it HeadIntensity) Get() api.Intensity {
	if it.channel != nil {
		return api.Intensity{it.channel.GetValue()}
	} else {
		return api.Intensity{}
	}
}

func (it HeadIntensity) Set(intensity api.Intensity) api.Intensity {
	return api.Intensity{it.channel.SetValue(intensity.Intensity)}
}

func (headIntensity *HeadIntensity) makeAPI() api.Intensity {
	return headIntensity.Get()
}

func (headIntensity HeadIntensity) GetREST() (web.Resource, error) {
	return headIntensity.makeAPI(), nil
}

func (headIntensity HeadIntensity) PostREST() (web.Resource, error) {
	return headIntensity.makeAPI(), nil
}

// Group.Intensity
type GroupIntensity struct {
	heads map[api.HeadID]HeadIntensity
}

func (groupIntensity GroupIntensity) exists() bool {
	return len(groupIntensity.heads) > 0
}

func (groupIntensity GroupIntensity) Get() (intensity api.Intensity) {
	for _, headIntensity := range groupIntensity.heads {
		return headIntensity.Get()
	}
	return
}

func (groupIntensity GroupIntensity) Set(intensity api.Intensity) api.Intensity {
	for _, headIntensity := range groupIntensity.heads {
		headIntensity.Set(intensity)
	}
	return intensity
}

func (groupIntensity *GroupIntensity) makeAPI() api.Intensity {
	return groupIntensity.Get()
}

// Web API
type APIIntensity struct {
	headIntensity  *HeadIntensity
	groupIntensity *GroupIntensity

	ScaleIntensity *api.Value
	api.Intensity
}

func (apiIntensity APIIntensity) IsZero() bool {
	return apiIntensity.Intensity == 0.0
}
func (apiIntensity APIIntensity) Equals(other APIIntensity) bool {
	return apiIntensity.Intensity == other.Intensity
}

func (apiIntensity *APIIntensity) initHead(headIntensity *HeadIntensity) error {
	if headIntensity == nil {
		return fmt.Errorf("Head does not support intensity")
	}

	apiIntensity.headIntensity = headIntensity

	return nil
}

func (apiIntensity *APIIntensity) initGroup(groupIntensity *GroupIntensity) error {
	if groupIntensity == nil {
		return web.RequestErrorf("Group does not support intensity")
	}

	apiIntensity.groupIntensity = groupIntensity

	return nil
}

func (apiIntensity *APIIntensity) Apply() error {
	if apiIntensity.ScaleIntensity != nil {
		apiIntensity.Intensity = apiIntensity.Intensity.ScaleIntensity(*apiIntensity.ScaleIntensity)
	}

	if apiIntensity.headIntensity != nil {
		apiIntensity.Intensity = apiIntensity.headIntensity.Set(apiIntensity.Intensity)
	}

	if apiIntensity.groupIntensity != nil {
		apiIntensity.Intensity = apiIntensity.groupIntensity.Set(apiIntensity.Intensity)
	}

	return nil
}

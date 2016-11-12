package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qmsk/go-web"
)

type Intensity Value // 0.0 .. 1.0

type HeadIntensity struct {
	channel *Channel
}

func (it HeadIntensity) exists() bool {
	return it.channel != nil
}

func (it HeadIntensity) Get() Intensity {
	if it.channel != nil {
		return Intensity(it.channel.GetValue())
	} else {
		return Intensity(INVALID)
	}
}

func (it HeadIntensity) Set(intensity Intensity) Intensity {
	return Intensity(it.channel.SetValue(Value(intensity)))
}

// Web API
type APIHeadIntensity struct {
	headIntensity *HeadIntensity

	Intensity
}

func (headIntensity *HeadIntensity) makeAPI() *APIHeadIntensity {
	if headIntensity == nil {
		return nil
	}

	return &APIHeadIntensity{
		headIntensity: headIntensity,
		Intensity:     headIntensity.Get(),
	}
}

func (headIntensity HeadIntensity) GetREST() (web.Resource, error) {
	return headIntensity.makeAPI(), nil
}

func (headIntensity HeadIntensity) PostREST() (web.Resource, error) {
	return headIntensity.makeAPI(), nil
}

func (apiHeadIntensity *APIHeadIntensity) Apply() error {
	if apiHeadIntensity.headIntensity == nil {
		return web.RequestErrorf("Head does not support intensity")
	}

	log.Debugln("heads:APIHeadIntensity.Apply", apiHeadIntensity.Intensity)

	apiHeadIntensity.Intensity = apiHeadIntensity.headIntensity.Set(apiHeadIntensity.Intensity)

	return nil
}

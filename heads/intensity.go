package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-web"
)

type Intensity Value // 0.0 .. 1.0

type HeadIntensity struct {
	channel *Channel
}

func (it HeadIntensity) Exists() bool {
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
	headIntensity HeadIntensity

	Intensity
}

func (headIntensity HeadIntensity) makeAPI() *APIHeadIntensity {
	if !headIntensity.Exists() {
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
	log.Debugln("heads:APIHeadIntensity.Apply", apiHeadIntensity.Intensity)

	apiHeadIntensity.Intensity = apiHeadIntensity.headIntensity.Set(apiHeadIntensity.Intensity)

	return nil
}

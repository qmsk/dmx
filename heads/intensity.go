package heads

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

func (it HeadIntensity) Set(intensity Intensity) {
	it.channel.SetValue(Value(intensity))
}

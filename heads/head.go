package heads

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
	"github.com/SpComb/qmsk-web"
)

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	id       string
	config   HeadConfig
	headType *HeadType
	output   *Output

	channels   HeadChannels
	parameters HeadParameters
}

type HeadChannels map[ChannelType]*Channel

type HeadParameters struct {
	Intensity *HeadIntensity `json:"intensity,omitempty"`
	Color     *HeadColor     `json:"color,omitempty"`
}

func (head *Head) Name() string {
	return head.config.Name
}

func (head *Head) String() string {
	return fmt.Sprintf("%v@%v[%d]", head.id, head.output, head.config.Address)
}

func (head *Head) init() {
	head.channels = make(HeadChannels)

	for channelIndex, channelType := range head.headType.Channels {
		var channel = &Channel{
			channelType: channelType,
			index:       uint(channelIndex),
			output:      head.output,
			address:     head.config.Address + dmx.Address(channelIndex),
		}

		channel.init()

		head.channels[channelType] = channel
	}

	// set parameters
	if headIntensity := head.getIntensity(); headIntensity.exists() {
		head.parameters.Intensity = &headIntensity
	}
	if headColor := head.getColor(); headColor.exists() {
		head.parameters.Color = &headColor
	}
}

func (head *Head) getChannel(channelType ChannelType) *Channel {
	return head.channels[channelType]
}

func (head *Head) getIntensity() HeadIntensity {
	return HeadIntensity{
		channel: head.getChannel(ChannelType{Intensity: true}),
	}
}

func (head *Head) getColor() HeadColor {
	return HeadColor{
		red:       head.getChannel(ChannelType{Color: ColorChannelRed}),
		green:     head.getChannel(ChannelType{Color: ColorChannelGreen}),
		blue:      head.getChannel(ChannelType{Color: ColorChannelBlue}),
		intensity: head.getChannel(ChannelType{Intensity: true}),
	}
}

func (head *Head) Parameters() HeadParameters {
	return head.parameters
}

// web API
type APIHeadParameters struct {
	head *Head

	Intensity *APIHeadIntensity `json:",omitempty"`
	Color     *APIHeadColor     `json:",omitempty"`
}

func (headParameters HeadParameters) makeAPI() APIHeadParameters {
	return APIHeadParameters{
		Intensity: headParameters.Intensity.makeAPI(),
		Color:     headParameters.Color.makeAPI(),
	}
}

func (headParameters HeadParameters) GetREST() (web.Resource, error) {
	return headParameters.makeAPI(), nil
}

type APIHead struct {
	ID     string
	Config HeadConfig
	Type   *HeadType

	Channels APIChannels
	APIHeadParameters
}

func (head *Head) makeAPIChannels() APIChannels {
	var apiChannels APIChannels

	for _, channelType := range head.headType.Channels {
		apiChannels = append(apiChannels, head.getChannel(channelType).makeAPI())
	}

	return apiChannels
}

func (head *Head) makeAPI() APIHead {
	log.Debugln("heads:Head.makeAPI", head)

	return APIHead{
		ID:     head.id,
		Config: head.config,
		Type:   head.headType,

		Channels:          head.makeAPIChannels(),
		APIHeadParameters: head.parameters.makeAPI(),
	}
}

func (head *Head) GetREST() (web.Resource, error) {
	return head.makeAPI(), nil
}

func (head *Head) PostREST() (web.Resource, error) {
	// parameters only, not configuration
	return &APIHeadParameters{head: head}, nil
}

func (params APIHeadParameters) Apply() error {
	log.Debugln("heads:Head.Apply", params.head,
		"intensity", params.Intensity,
		"color", params.Color,
	)

	if params.Intensity != nil {
		params.Intensity.headIntensity = params.head.parameters.Intensity

		if err := params.Intensity.Apply(); err != nil {
			return err
		}
	}

	if params.Color != nil {
		params.Color.headColor = params.head.parameters.Color

		if err := params.Color.Apply(); err != nil {
			return err
		}
	}

	return nil
}

func (head *Head) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return head.parameters, nil
	case "channels":
		return head.channels, nil
	case "intensity":
		return head.parameters.Intensity, nil
	case "color":
		return head.parameters.Color, nil
	default:
		return nil, nil
	}
}

func (headChannels HeadChannels) makeAPI() APIChannels {
	var apiChannels = make(APIChannels, len(headChannels))

	for _, channel := range headChannels {
		apiChannels[channel.index] = channel.makeAPI()
	}

	return apiChannels
}

func (headChannels HeadChannels) GetREST() (web.Resource, error) {
	log.Debugln("heads:HeadChannels.GetREST")

	return headChannels.makeAPI(), nil
}

func (headChannels HeadChannels) Index(name string) (web.Resource, error) {
	for channelType, channel := range headChannels {
		if channelType.String() == name {
			log.Debugln("heads:HeadChannels.Index", name, channel)

			return web.GetPostResource(channel), nil
		}
	}

	log.Debugln("heads:HeadChannels.Index", name, nil)

	return nil, nil
}

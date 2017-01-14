package heads

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
	"github.com/qmsk/go-web"
)

type HeadChannels map[ChannelType]*Channel

func (headChannels HeadChannels) GetID(id string) *Channel {
	for channelType, channel := range headChannels {
		if channelType.String() == id {
			return channel
		}
	}

	return nil
}

func (headChannels HeadChannels) makeAPI() APIChannels {
	var apiChannels = make(APIChannels)

	for channelType, channel := range headChannels {
		apiChannels[channelType.String()] = channel.makeAPI()
	}

	return apiChannels
}

func (headChannels HeadChannels) GetREST() (web.Resource, error) {
	log.Debugln("heads:HeadChannels.GetREST")

	return headChannels.makeAPI(), nil
}

func (headChannels HeadChannels) Index(name string) (web.Resource, error) {
	if channel := headChannels.GetID(name); channel == nil {
		return nil, nil
	} else {
		return web.GetPostResource(channel), nil
	}
}

type HeadParameters struct {
	Intensity *HeadIntensity `json:"intensity,omitempty"`
	Color     *HeadColor     `json:"color,omitempty"`
}

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	id       HeadID
	config   HeadConfig
	headType *HeadType
	output   *Output
	events   *Events

	channels   HeadChannels
	parameters HeadParameters
}

func (head *Head) String() string {
	return fmt.Sprintf("%v", head.id)
}

func (head *Head) Name() string {
	return head.config.Name
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

// Web API GET
type APIHead struct {
	ID     HeadID
	Config HeadConfig
	Type   *HeadType

	Channels  map[string]APIChannel `json:",omitempty"`
	Intensity *APIHeadIntensity     `json:",omitempty"`
	Color     *APIHeadColor         `json:",omitempty"`
}

func (head *Head) makeAPI() APIHead {
	log.Debugln("heads:Head.makeAPI", head)

	return APIHead{
		ID:     head.id,
		Config: head.config,
		Type:   head.headType,

		Channels:  head.channels.makeAPI(),
		Intensity: head.parameters.Intensity.makeAPI(),
		Color:     head.parameters.Color.makeAPI(),
	}
}

func (head *Head) GetREST() (web.Resource, error) {
	return head.makeAPI(), nil
}

// Web API POST
type APIHeadParams struct {
	head *Head

	Channels  map[string]APIChannelParams `json:",omitempty"`
	Intensity *APIHeadIntensity           `json:",omitempty"`
	Color     *APIHeadColor               `json:",omitempty"`
}

func (head *Head) PostREST() (web.Resource, error) {
	// parameters only, not configuration
	return &APIHeadParams{head: head}, nil
}

func (post *APIHeadParams) Apply() error {
	log.Debugln("heads:Head.Apply", post.head,
		"channels", post.Channels,
		"intensity", post.Intensity,
		"color", post.Color,
	)

	for channelID, channelParams := range post.Channels {
		if channel := post.head.channels.GetID(channelID); channel == nil {
			return web.Errorf(404, "Channel not found: %v", channelID)
		} else {
			channelParams.channel = channel
		}

		if err := channelParams.Apply(); err != nil {
			return err
		}
	}

	if post.Intensity != nil {
		post.Intensity.headIntensity = post.head.parameters.Intensity

		if err := post.Intensity.Apply(); err != nil {
			return err
		}
	}

	if post.Color != nil {
		post.Color.headColor = post.head.parameters.Color

		if err := post.Color.Apply(); err != nil {
			return err
		}
	}

	return nil
}

func (head *Head) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return head, nil
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

// Web API Events
func (head *Head) Apply() error {
	head.events.updateHead(head.String(), head.makeAPI())

	return nil
}

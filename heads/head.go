package heads

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
	"github.com/SpComb/qmsk-web"
)

type Channel struct {
	channelType ChannelType
	output      *Output

	address dmx.Address
}

func (channel *Channel) init() {
	channel.output.SetDMX(channel.address, 0)
}

func (channel *Channel) GetDMX() dmx.Channel {
	return channel.output.GetDMX(channel.address)
}
func (channel *Channel) GetValue() Value {
	return channel.output.GetValue(channel.address)
}

func (channel *Channel) SetDMX(value dmx.Channel) {
	channel.output.SetDMX(channel.address, value)
}
func (channel *Channel) SetValue(value Value) Value {
	return channel.output.SetValue(channel.address, value)
}

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	id       string
	config   HeadConfig
	headType *HeadType
	output   *Output

	channels map[ChannelType]*Channel
}

func (head *Head) Name() string {
	return head.config.Name
}

func (head *Head) String() string {
	return fmt.Sprintf("%v@%v[%d]", head.id, head.output, head.config.Address)
}

func (head *Head) init() {
	head.channels = make(map[ChannelType]*Channel)

	for channelOffset, channelType := range head.headType.Channels {
		var channel = &Channel{
			channelType: channelType,
			output:      head.output,
			address:     head.config.Address + dmx.Address(channelOffset),
		}

		channel.init()

		head.channels[channelType] = channel
	}
}

func (head *Head) getChannel(channelType ChannelType) *Channel {
	return head.channels[channelType]
}

func (head *Head) Intensity() HeadIntensity {
	return HeadIntensity{
		channel: head.getChannel(ChannelType{Intensity: true}),
	}
}

func (head *Head) Color() HeadColor {
	return HeadColor{
		red:       head.getChannel(ChannelType{Color: ColorChannelRed}),
		green:     head.getChannel(ChannelType{Color: ColorChannelGreen}),
		blue:      head.getChannel(ChannelType{Color: ColorChannelBlue}),
		intensity: head.getChannel(ChannelType{Intensity: true}),
	}
}

// web API
type APIHead struct {
	ID     string
	Config HeadConfig
	Type   *HeadType

	Intensity *APIHeadIntensity `json:"intensity,omitempty"`
	Color     *APIHeadColor     `json:"color,omitempty"`
}

func (head *Head) makeAPI() APIHead {
	log.Debugln("heads:Head.makeAPI", head)

	return APIHead{
		ID:     head.id,
		Config: head.config,
		Type:   head.headType,

		Intensity: head.Intensity().makeAPI(),
		Color:     head.Color().makeAPI(),
	}
}

func (head *Head) GetREST() (web.Resource, error) {
	return head.makeAPI(), nil
}

func (head *Head) Index(name string) (web.Resource, error) {
	switch name {
	case "intensity":
		return head.Intensity(), nil
	case "color":
		return head.Color(), nil
	default:
		return nil, nil
	}
}

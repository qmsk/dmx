package heads

import (
	"github.com/SpComb/qmsk-dmx"
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
func (channel *Channel) SetValue(value Value) {
	channel.output.SetValue(channel.address, value)
}

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	headType *HeadType
	address  dmx.Address

	channels map[ChannelType]*Channel
}

func (head *Head) init(output *Output, headType *HeadType) {
	head.channels = make(map[ChannelType]*Channel)

	for channelOffset, channelType := range headType.Channels {
		var channel = &Channel{
			channelType: channelType,
			output:      output,
			address:     head.address + dmx.Address(channelOffset),
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

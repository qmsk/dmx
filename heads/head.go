package heads

import (
	"github.com/SpComb/qmsk-dmx"
)

type Channel struct {
	config ChannelConfig
	output *Output

	address dmx.Address
}

func (channel *Channel) init() {
	channel.output.Set(channel.address, 0)
}

func (channel *Channel) Get() dmx.Channel {
	return channel.output.Get(channel.address)
}

func (channel *Channel) Set(value dmx.Channel) {
	channel.output.Set(channel.address, value)
}

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	headType *HeadType
	address  dmx.Address

	channels map[ChannelConfig]*Channel
}

func (head *Head) init(output *Output, headType *HeadType) {
	head.channels = make(map[ChannelConfig]*Channel)

	for channelOffset, channelConfig := range headType.Channels {
		var channel = &Channel{
			config:  channelConfig,
			output:  output,
			address: head.address + dmx.Address(channelOffset),
		}

		channel.init()

		head.channels[channelConfig] = channel
	}
}

func (head *Head) IntensityChannel() *Channel {
	return head.channels[ChannelConfig{Class: ChannelClassIntensity}]
}

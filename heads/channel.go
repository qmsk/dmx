package heads

import "github.com/SpComb/qmsk-dmx"

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

// Web API
type APIChannel struct {
	channel *Channel

	Type    ChannelType
	Address dmx.Address

	DMX   dmx.Channel
	Value Value
}

func (channel *Channel) makeAPI() APIChannel {
	return APIChannel{
		channel: channel,
		Type:    channel.channelType,
		Address: channel.address,
		DMX:     channel.GetDMX(),
		Value:   channel.GetValue(),
	}
}

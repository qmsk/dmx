package heads

import (
	"github.com/qmsk/dmx"
	"github.com/qmsk/go-web"
)

type ChannelType struct {
	Control   string       `json:",omitempty"`
	Intensity bool         `json:",omitempty"`
	Color     ColorChannel `json:",omitempty"`
}

func (channelType ChannelType) String() string {
	if channelType.Control != "" {
		return "control:" + channelType.Control
	}
	if channelType.Intensity {
		return "intensity"
	}
	if channelType.Color != "" {
		return "color:" + string(channelType.Color)
	}

	return ""
}

type Channel struct {
	channelType ChannelType
	output      *Output
	index       uint

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

	ID      string
	Type    ChannelType
	Index   uint
	Address dmx.Address

	DMX   dmx.Channel
	Value Value
}

func (channel *Channel) makeAPI() APIChannel {
	return APIChannel{
		channel: channel,
		ID:      channel.channelType.String(),
		Index:   channel.index,
		Type:    channel.channelType,
		Address: channel.address,
		DMX:     channel.GetDMX(),
		Value:   channel.GetValue(),
	}
}

func (channel *Channel) GetREST() (web.Resource, error) {
	return channel.makeAPI(), nil
}

type APIChannels map[string]APIChannel

func (apiChannels APIChannels) GetREST() (web.Resource, error) {
	return apiChannels, nil
}

type APIChannelParams struct {
	channel *Channel

	DMX   *dmx.Channel `json:",omitempty"`
	Value *Value       `json:",omitempty"`
}

func (channel *Channel) PostREST() (web.Resource, error) {
	return &APIChannelParams{channel: channel}, nil
}

func (params *APIChannelParams) Apply() error {
	if params.DMX != nil {
		params.channel.SetDMX(*params.DMX)
	}
	if params.Value != nil {
		*params.Value = params.channel.SetValue(*params.Value)
	}

	if params.DMX == nil {
		var dmxValue = params.channel.GetDMX()
		params.DMX = &dmxValue
	}
	if params.Value == nil {
		var value = params.channel.GetValue()
		params.Value = &value
	}

	return nil
}

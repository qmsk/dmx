package heads

import (
	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/go-web"
)

// Channels
type channels map[api.ChannelID]*Channel

func (channels channels) makeAPI() api.Channels {
	var apiChannels = make(api.Channels)

	for id, channel := range channels {
		apiChannels[id] = channel.makeAPI()
	}

	return apiChannels
}

func (channels channels) GetREST() (web.Resource, error) {
	return channels.makeAPI(), nil
}

func (channels channels) Index(name string) (web.Resource, error) {
	if channel := channels[api.ChannelID(name)]; channel == nil {
		return nil, nil
	} else {
		return channel, nil
	}
}

type Channel struct {
	id     api.ChannelID
	config api.ChannelConfig
	output *Output
	index  uint

	address dmx.Address
}

func (channel *Channel) init() {
	channel.output.SetDMX(channel.address, 0)
}

func (channel *Channel) GetDMX() api.DMXValue {
	return api.DMXValue(channel.output.GetDMX(channel.address))
}
func (channel *Channel) GetValue() api.Value {
	return api.Value(channel.output.GetValue(channel.address))
}

func (channel *Channel) SetDMX(value api.DMXValue) {
	channel.output.SetDMX(channel.address, dmx.Channel(value))
}
func (channel *Channel) SetValue(value api.Value) api.Value {
	return api.Value(channel.output.SetValue(channel.address, Value(value)))
}

func (channel *Channel) makeAPI() api.Channel {
	return api.Channel{
		ID:      channel.id,
		Index:   channel.index,
		Config:  channel.config,
		Address: api.DMXAddress(channel.address),
		DMX:     channel.GetDMX(),
		Value:   channel.GetValue(),
	}
}

func (channel *Channel) GetREST() (web.Resource, error) {
	return channel.makeAPI(), nil
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

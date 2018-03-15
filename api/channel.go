package api

type ChannelColor string

const (
	ChannelColorRed   = "red"
	ChannelColorGreen = "green"
	ChannelColorBlue  = "blue"
)

type ChannelConfig struct {
	Control   string       `json:",omitempty"`
	Intensity bool         `json:",omitempty"`
	Color     ChannelColor `json:",omitempty"`
}

type ChannelID string
type Channels map[ChannelID]Channel

type Channel struct {
	ID      ChannelID
	Config  ChannelConfig
	Index   uint
	Address DMXAddress

	DMX   DMXValue
	Value Value
}

type ChannelParams struct {
	DMX   *DMXValue `json:",omitempty"`
	Value *Value    `json:",omitempty"`
}

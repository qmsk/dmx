package api

type TypeID string
type HeadType struct {
	Vendor string
	Model  string
	Mode   string
	URL    string

	Channels []ChannelConfig
	Colors   Colors
}

type HeadConfig struct {
	Type     TypeID
	Universe Universe
	Address  DMXAddress
	Name     string
	Count    uint // Clone multiple copies of the head at id.N
	Groups   []GroupID
}

type HeadID string

type Heads map[HeadID]Head

type Head struct {
	ID     HeadID
	Config HeadConfig
	Type   HeadType

	Channels  map[string]Channel `json:",omitempty"`
	Intensity *Intensity         `json:",omitempty"`
	Color     *Color             `json:",omitempty"`
}

type HeadParams struct {
	Channels  map[string]ChannelParams `json:",omitempty"`
	Intensity *Intensity               `json:",omitempty"`
	Color     *Color                   `json:",omitempty"`
}

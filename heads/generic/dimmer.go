package generic

import "github.com/SpComb/qmsk-dmx/heads"

var Dimmer = heads.HeadConfig{
	Vendor: "Generic",
	Model:  "Dimmer",

	Channels: []heads.ChannelConfig{
		heads.ChannelConfig{
			Class: heads.ChannelClassIntensity,
		},
	},
}

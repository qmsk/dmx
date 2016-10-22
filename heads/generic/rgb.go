package generic

import "github.com/SpComb/qmsk-dmx/heads"

var RGB = heads.HeadConfig{
	Vendor: "Generic",
	Model:  "RGB",

	Channels: []heads.ChannelConfig{
		heads.ChannelConfig{
			Class: heads.ChannelClassColor,
			Name:  "red",
		},
		heads.ChannelConfig{
			Class: heads.ChannelClassColor,
			Name:  "green",
		},
		heads.ChannelConfig{
			Class: heads.ChannelClassColor,
			Name:  "blue",
		},
	},
}

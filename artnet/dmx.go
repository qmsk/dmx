package artnet

import (
	"net"

	dmx "github.com/SpComb/qmsk-dmx"
)

type DMXUniverse []uint8

type ArtDmx struct {
	ArtHeader
	ProtVer ProtVer

	Sequence uint8
	Physical uint8
	SubUni   uint8
	Net      uint8
	Length   uint16
}

func (transport *Transport) SendDMX(addr *net.UDPAddr, sequence uint8, address Address, universe dmx.Universe) error {
	packet := ArtDmx{
		ArtHeader: ArtHeader{
			ID:     ARTNET,
			OpCode: OpDmx,
		},
		ProtVer:  ProtVer14,
		Sequence: sequence,
		SubUni:   address.SubUni,
		Net:      address.Net,
		Length:   uint16(len(universe)),
	}

	return transport.send(addr, packet, universe.Bytes())
}

package artnet

import (
	"net"
)

type Universe []uint8

type ArtDmx struct {
	ArtHeader
	ProtVer uint16

	Sequence uint8
	Physical uint8
	SubUni   uint8
	Net      uint8
	Length   uint16
}

func (transport *Transport) SendDMX(addr *net.UDPAddr, sequence uint8, address Address, data Universe) error {
	packet := ArtDmx{
		ArtHeader: ArtHeader{
			ID:     ARTNET,
			OpCode: OpDmx,
		},
		ProtVer:  ProtVer,
		Sequence: sequence,
		SubUni:   address.SubUni,
		Net:      address.Net,
		Length:   uint16(len(data)),
	}

	return transport.send(addr, packet, data)
}

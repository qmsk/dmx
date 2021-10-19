package artnet

import (
	"net"
)

type ArtSync struct {
	ArtHeader
	ProtVer ProtVer

	Aux1 uint8
	Aux2 uint8
}

func (transport *Transport) SendSync(addr *net.UDPAddr) error {
	packet := ArtSync{
		ArtHeader: ArtHeader{
			ID:     ARTNET,
			OpCode: OpSync,
		},
		ProtVer: ProtVer14,
	}

	return transport.send(addr, packet, nil)
}

package artnet

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net"
)

type ArtPoll struct {
	ArtHeader
	ProtVer ProtVer

	TalkToMe uint8
	Priority uint8
}

func (transport *Transport) SendPoll(addr *net.UDPAddr) error {
	return transport.send(addr, ArtPoll{
		ArtHeader: ArtHeader{
			ID:     ARTNET,
			OpCode: OpPoll,
		},
		ProtVer: ProtVer14,
	}, nil)
}

type ArtPollReply struct {
	ArtHeader

	IPAddress   [4]byte
	PortNumber  uint16 // XXX: swapped byte order
	VersInfo    uint16
	NetSwitch   uint8
	SubSwitch   uint8
	Oem         uint16
	UbeaVersion uint8
	Status1     uint8
	EstaMan     uint16
	ShortName   [18]byte
	LongName    [64]byte
	NodeReport  [64]byte
	NumPorts    uint16
	PortTypes   [4]uint8
	GoodInput   [4]uint8
	GoodOutput  [4]uint8
	SwIn        [4]uint8
	SwOut       [4]uint8
	SwVideo     uint8
	SwMacro     uint8
	SwRemote    uint8
	Spare1      byte
	Spare2      byte
	Spare3      byte
	Style       byte
	Mac         [6]byte
	BindIp      [4]byte
	BindIndex   uint8
	Status2     uint8
}

func decodeMac(mac [6]byte) string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func (p ArtPollReply) NodeConfig() NodeConfig {
	var address = Address{
		Net:    p.NetSwitch,
		SubUni: p.SubSwitch,
	}

	var nodeConfig = NodeConfig{
		OEM:         p.Oem,
		Version:     p.VersInfo,
		Name:        decodeString(p.ShortName[:]),
		Description: decodeString(p.LongName[:]),
		Report:      decodeString(p.NodeReport[:]),
		Ethernet:    decodeMac(p.Mac),
		BaseAddress: address,
	}

	for i := 0; i < int(p.NumPorts) && i < 4; i++ {
		log.Debugf("decode poll reply: port=%d type=%04x", i, p.PortTypes[i])

		if p.PortTypes[i]&0x80 != 0 {
			nodeConfig.OutputPorts = append(nodeConfig.OutputPorts, OutputPort{
				Address: Address{
					Net:    address.Net,
					SubUni: address.SubUni | p.SwOut[i],
				},
				Type:   p.PortTypes[i] & 0x1F,
				Status: p.GoodOutput[i],
			})
		}
		if p.PortTypes[i]&0x40 != 0 {
			nodeConfig.InputPorts = append(nodeConfig.InputPorts, InputPort{
				Address: Address{
					Net:    address.Net,
					SubUni: address.SubUni | p.SwIn[i],
				},
				Type:   p.PortTypes[i] & 0x1F,
				Status: p.GoodInput[i],
			})
		}
	}

	return nodeConfig
}

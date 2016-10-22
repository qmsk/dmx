package artnet

import (
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	dmx "github.com/SpComb/qmsk-dmx"
)

type InputPort struct {
	Address Address

	Type   uint8
	Status uint8
}

type OutputPort struct {
	Address Address

	Type   uint8
	Status uint8
}

type NodeConfig struct {
	OEM         uint16
	Version     uint16
	Name        string
	Description string
	Report      string
	Ethernet    string

	BaseAddress Address
	InputPorts  []InputPort
	OutputPorts []OutputPort
}

type Node struct {
	log *log.Entry

	timeout time.Duration

	transport *Transport
	addr      *net.UDPAddr // unicast

	config        NodeConfig
	sequence      uint8
	discoveryTime time.Time
}

func (node *Node) String() string {
	return node.addr.String()
}

func (node *Node) Config() NodeConfig {
	// XXX: atomic
	return node.config
}

func (node *Node) SendDMX(address Address, universe dmx.Universe) error {
	node.sequence++

	if node.sequence == 0 {
		node.sequence = 1
	}

	node.log.Debugf("SendDMX %v @ %v", address, node.sequence)

	return node.transport.SendDMX(node.addr, node.sequence, address, universe)
}

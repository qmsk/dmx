package artnet

import (
	"net"
	"time"

	dmx "github.com/qmsk/dmx"
	"github.com/qmsk/dmx/logging"
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

func (controller *Controller) makeNode(addr *net.UDPAddr, config NodeConfig) (*Node, error) {
	var node = Node{
		log:       controller.config.Log.Logger("node", addr.String()),
		timeout:   controller.config.DiscoveryTimeout,
		transport: controller.transport,
		addr:      addr,
		config:    config,
	}

	return &node, nil
}

type Node struct {
	log logging.Logger

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
	// start sequence at 0
	var sequence = node.sequence

	node.sequence++

	if node.sequence == 0 {
		node.sequence = 1
	}

	node.log.Debugf("SendDMX %v @ %v [%d]", address, sequence, len(universe))

	return node.transport.SendDMX(node.addr, sequence, address, universe)
}

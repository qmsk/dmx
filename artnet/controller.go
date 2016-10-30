package artnet

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
	dmx "github.com/SpComb/qmsk-dmx"
)

type Config struct {
	Listen    string   `long:"artnet-listen" value-name:"ADDR" default:"0.0.0.0"`
	Discovery []string `long:"artnet-discovery" value-name:"ADDR" default:"255.255.255.255"`

	DiscoveryInterval time.Duration `long:"artnet-discovery-interval" value-name:"DURATION" default:"3s"`
	DiscoveryTimeout  time.Duration `long:"artnet-discovery-timeout" value-name:"DURATION" default:"3s"`

	DMXRefresh time.Duration `long:"artnet-dmx-refresh" value-name:"DURATION" default:"1s"`
}

func (config Config) Controller() (*Controller, error) {
	var controller = Controller{
		log: log.WithFields(log.Fields{"prefix": "artnet:Controller"}),

		config: config,

		universes: make(map[Address]*Universe),
	}

	listenAddr := net.JoinHostPort(config.Listen, fmt.Sprintf("%d", Port))

	if udpAddr, err := net.ResolveUDPAddr("udp", listenAddr); err != nil {
		return nil, err
	} else if udpConn, err := net.ListenUDP("udp", udpAddr); err != nil {
		return nil, err
	} else {
		controller.transport = &Transport{
			udpConn: udpConn,
		}
	}

	for _, discovery := range config.Discovery {
		discoveryAddr := net.JoinHostPort(discovery, fmt.Sprintf("%d", Port))

		if udpAddr, err := net.ResolveUDPAddr("udp", discoveryAddr); err != nil {
			return nil, err
		} else {
			controller.discoveryAddrs = append(controller.discoveryAddrs, udpAddr)
		}
	}

	return &controller, nil
}

type pollEvent struct {
	recvTime time.Time
	srcAddr  *net.UDPAddr
	packet   ArtPollReply
}

func (event pollEvent) String() string {
	return event.srcAddr.String()
}

type Controller struct {
	log *log.Entry

	config Config

	transport *Transport
	pollChan  chan pollEvent

	// discovery handler
	discoveryAddrs []*net.UDPAddr // sending to unicast/broadcast addresses
	discoveryState atomic.Value
	discoveryChan  chan Discovery

	// state
	universes map[Address]*Universe
}

func (controller *Controller) Start(discoveryChan chan Discovery) {
	controller.pollChan = make(chan pollEvent)
	controller.discoveryChan = discoveryChan

	go controller.recv()
	go controller.discovery(controller.pollChan)
}

func (controller *Controller) recv() {
	for {
		if packet, srcAddr, err := controller.transport.recv(); err != nil {
			// XXX: fatal if socket is dead?
			controller.log.Errorf("recv %v: %v", srcAddr, err)
		} else if err := controller.recvPacket(packet, srcAddr); err != nil {
			controller.log.Warnf("recv %v: %v", srcAddr, err)
		}
	}
}

func (controller *Controller) recvPacket(packet ArtPacket, srcAddr *net.UDPAddr) error {
	switch packetType := packet.(type) {
	case *ArtPoll:
		if !ProtVer14.IsCompatible(packetType.ProtVer) {
			return fmt.Errorf("Invalid protocol version: %v < %v", packetType.ProtVer, ProtVer14)
		}

		// ignore
		return nil

	case *ArtPollReply:
		if controller.pollChan != nil {
			controller.pollChan <- pollEvent{
				recvTime: time.Now(),
				srcAddr:  srcAddr,
				packet:   *packetType,
			}
		}

	case *ArtDmx:
		if !ProtVer14.IsCompatible(packetType.ProtVer) {
			return fmt.Errorf("Invalid protocol version: %v < %v", packetType.ProtVer, ProtVer14)
		}

		// ignore
		return nil
	}

	return nil
}

// Send DMX universe using either unicast or broadcast.
//
// Implements ArtNet Subscription using the discovery nodes.
//
// If we have discovered a Node configured for the given address, the DMX packet is unicast to each such node.
// Otherwise, the packet is broadcast to the discovery address.
func (controller *Controller) SendDMX(address Address, universe dmx.Universe) error {
	discovery := controller.Discovery()

	var matchNodes = false

	for _, node := range discovery.Nodes {
		var matchNode = false

		for _, outputPort := range node.config.OutputPorts {
			if outputPort.Address == address {
				matchNode = true
			}
		}

		if matchNode {
			matchNodes = true

			// send unicast to node; may have multiple outputs for the same universe
			if err := node.SendDMX(address, universe); err != nil {
				return fmt.Errorf("Node %v: SendDMX %v: %v", node, address, err)
			}
		}
	}

	if !matchNodes {
		// send broadcast, did not find specific node
		for _, addr := range controller.discoveryAddrs {
			if err := controller.transport.SendDMX(addr, 0, address, universe); err != nil {
				return fmt.Errorf("SendDMX broadcast %v: %v", address, err)
			} else {
				controller.log.Debugf("SendDMX %v [%v]: broadcast %v ", address, len(universe), addr)
			}
		}
	}

	return nil
}

package artnet

import (
  "fmt"
  "net"
  log "github.com/Sirupsen/logrus"
  "time"
)

type Config struct {
  Listen    string    `long:"artnet-listen" value-name:"ADDR" default:"0.0.0.0"`
  Discovery string    `long:"artnet-discovery" value-name:"ADDR" default:"255.255.255.255"`

  DiscoveryInterval time.Duration `long:"artnet-discovery-interval" value-name:"DURATION" default:"3s"`
  DiscoveryTimeout  time.Duration `long:"artnet-discovery-timeout" value-name:"DURATION" default:"3s"`
}

func (config Config) Controller() (*Controller, error) {
  var controller = Controller{
    log:  log.WithFields(log.Fields{"prefix": "artnet:Controller"}),

    discoveryInterval: config.DiscoveryInterval,
  }

  listenAddr  := net.JoinHostPort(config.Listen, fmt.Sprintf("%d", Port))
  discoveryAddr := net.JoinHostPort(config.Discovery, fmt.Sprintf("%d", Port))

  if udpAddr, err := net.ResolveUDPAddr("udp", listenAddr); err != nil {
    return nil, err
  } else if udpConn, err := net.ListenUDP("udp", udpAddr); err != nil {
    return nil, err
  } else {
    controller.transport = &Transport{
      udpConn: udpConn,
    }
  }

  if udpAddr, err := net.ResolveUDPAddr("udp", discoveryAddr); err != nil {
    return nil, err
  } else {
    controller.discoveryAddr = udpAddr
  }

  return &controller, nil
}

type Controller struct {
  log *log.Entry

  transport *Transport

  // discovery handler
  discoveryAddr    *net.UDPAddr    // sending to broadcast
  discoveryChan     chan discoveryEvent
  discoveryInterval time.Duration
}
func (controller *Controller) Run() {

  controller.discoveryChan = make(chan discoveryEvent)

  go controller.discovery(controller.discoveryChan)

  controller.recv()
}

func (controller *Controller) Start() {
  controller.discoveryChan = make(chan discoveryEvent)

  go controller.recv()
  go controller.discovery(controller.discoveryChan)
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
    if packetType.ProtVer < ProtVer {
      return fmt.Errorf("Invalid protocol version: %v < %v", packetType.ProtVer, ProtVer)
    }

    // ignore
    return nil

  case *ArtPollReply:
    if controller.discoveryChan != nil {
      controller.discoveryChan <- discoveryEvent{
        recvTime: time.Now(),
        srcAddr: srcAddr,
        pollReply: *packetType,
      }
    }

  case *ArtDmx:
    if packetType.ProtVer < ProtVer {
      return fmt.Errorf("Invalid protocol version: %v < %v", packetType.ProtVer, ProtVer)
    }

    // ignore
    return nil
  }

  return nil
}

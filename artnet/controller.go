package artnet

import (
  "sync/atomic"
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

    config: config,
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

type pollEvent struct {
  recvTime    time.Time
  srcAddr     *net.UDPAddr
  packet      ArtPollReply
}

func (event pollEvent) String() string {
  return event.srcAddr.String()
}

type Controller struct {
  log *log.Entry

  config    Config

  transport *Transport
  pollChan   chan pollEvent

  // discovery handler
  discoveryAddr    *net.UDPAddr    // sending to broadcast
  discoveryState    atomic.Value
  discoveryChan     chan Discovery
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
    if packetType.ProtVer < ProtVer {
      return fmt.Errorf("Invalid protocol version: %v < %v", packetType.ProtVer, ProtVer)
    }

    // ignore
    return nil

  case *ArtPollReply:
    if controller.pollChan != nil {
      controller.pollChan <- pollEvent{
        recvTime: time.Now(),
        srcAddr: srcAddr,
        packet: *packetType,
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

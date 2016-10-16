package artnet

import (
  "net"
  "time"
  log "github.com/Sirupsen/logrus"
)

type InputPort struct {
  Address     Address

  Type        uint8
  Status      uint8
}

type OutputPort struct {
  Address     Address

  Type        uint8
  Status      uint8
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
  addr *net.UDPAddr    // unicast

  config NodeConfig
  sequence uint8
  discoveryTime time.Time
}

func (node *Node) String() string {
  return node.addr.String()
}

func (node *Node) Config() NodeConfig {
  return node.config
}

func (node *Node) SendDMX(address Address, data Universe) error {
  node.sequence++

  if node.sequence == 0 {
    node.sequence = 1
  }

  return node.transport.send(node.addr, ArtDmx{
    ArtHeader: ArtHeader{
      ID: ARTNET,
      OpCode: OpDmx,
    },
    ProtVer: ProtVer,
    Sequence: node.sequence,
    SubUni:   address.SubUni,
    Net:      address.Net,
    Length:   uint16(len(data)),
    Data:     data,
  })
}

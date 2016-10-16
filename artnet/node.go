package artnet

import (
  "net"
  "time"
  log "github.com/Sirupsen/logrus"
)

type Node struct {
  log *log.Entry

  transport *Transport
  addr *net.UDPAddr    // unicast
  sequence uint8

  discoveryTime time.Time
  timeout time.Duration
}

func (node *Node) String() string {
  return node.addr.String()
}

func (node *Node) SendDMX(address UniverseAddress, data Universe) error {
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

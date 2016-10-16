package artnet

import (
  "fmt"
)

const Port = 6454
const MTU = 1500

var ARTNET = [8]byte{'A', 'R', 'T', '-', 'N', 'E', 'T', 0}

const ProtVer = 14

type Address struct {
  Net     uint8 // 0-128
  SubUni  uint8
}

func (a Address) String() string {
  return fmt.Sprintf("%d:%d.%d",
    a.Net,
    (a.SubUni >> 4),
    a.SubUni & 0x0F,
  )
}

type ArtHeader struct {
  ID      [8]byte
  OpCode  uint16
}

func (h ArtHeader) Header() ArtHeader {
  return h
}

type ArtPacket interface {
  Header() ArtHeader
}

const (
  opMask uint16 = 0xffff
  OpPoll        = 0x2000
  OpPollReply   = 0x2100
  OpDmx         = 0x5000
)

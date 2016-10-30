package artnet

const Port = 6454
const MTU = 1500

var ARTNET = [8]byte{'A', 'r', 't', '-', 'N', 'e', 't', 0}

const ProtVer = 14

type ArtHeader struct {
	ID     [8]byte
	OpCode uint16
}

func (h ArtHeader) Header() ArtHeader {
	return h
}

type ArtPacket interface {
	Header() ArtHeader
}

const (
	opMask      uint16 = 0xffff
	OpPoll             = 0x2000
	OpPollReply        = 0x2100
	OpDmx              = 0x5000
)

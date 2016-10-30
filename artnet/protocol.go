package artnet

const Port = 6454
const MTU = 1500

var ARTNET = [8]byte{'A', 'r', 't', '-', 'N', 'e', 't', 0}

type OpCode struct {
	Lo uint8
	Hi uint8
}
type ProtVer struct {
	Hi uint8
	Lo uint8
}

func (protVer ProtVer) ToUint() uint {
	return uint(protVer.Hi<<8) + uint(protVer.Lo)
}

func (protVer ProtVer) IsCompatible(otherVersion ProtVer) bool {
	return protVer.ToUint() == otherVersion.ToUint()
}

type ArtHeader struct {
	ID     [8]byte
	OpCode OpCode
}

func (h ArtHeader) Header() ArtHeader {
	return h
}

type ArtPacket interface {
	Header() ArtHeader
}

var ProtVer14 = ProtVer{0, 14}
var (
	OpPoll      = OpCode{0x00, 0x20}
	OpPollReply = OpCode{0x00, 0x21}
	OpDmx       = OpCode{0x00, 0x50}
)

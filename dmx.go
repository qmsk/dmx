package dmx

import (
	"encoding/hex"
)

type Address int
type Channel uint8

const UniverseChannels Address = 512

type Universe []Channel

func MakeUniverse() Universe {
	var universe Universe

	universe.init()

	return universe
}

func (universe *Universe) init() {
	*universe = make([]Channel, 0, UniverseChannels)
}

func (universe Universe) Bytes() []byte {
	var buf = make([]byte, len(universe))
	for i, channel := range universe {
		buf[i] = byte(channel)
	}
	return buf
}

func (universe Universe) String() string {
	return hex.Dump(universe.Bytes())
}

func (universe Universe) Copy() Universe {
	var out = make(Universe, len(universe))

	copy(out, universe)

	return out
}

func (universe Universe) Get(address Address) Channel {
	if address <= 0 || address > UniverseChannels {
		panic("Invalid DMX address")
	} else if int(address) > len(universe) {
		return 0
	}

	return universe[address-1]
}

func (universe *Universe) Set(address Address, value Channel) {
	if address <= 0 || address > UniverseChannels {
		panic("Invalid DMX address")
	} else if int(address) > len(*universe) {
		*universe = (*universe)[0:address]
	}

	(*universe)[address-1] = value
}

type Writer interface {
	WriteDMX(dmx Universe) error
}

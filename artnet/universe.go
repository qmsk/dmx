package artnet

import (
	"fmt"

	dmx "github.com/SpComb/qmsk-dmx"
)

type Address struct {
	Net    uint8 // 0-128
	SubUni uint8
}

func (a Address) String() string {
	return fmt.Sprintf("%d:%d.%d",
		a.Net,
		(a.SubUni >> 4),
		a.SubUni&0x0F,
	)
}

func (a Address) Integer() int {
	return int(uint(a.Net<<8) | uint(a.SubUni))
}

func (controller *Controller) Universes() map[Address]Universe {
	discovery := controller.Discovery()
	universes := make(map[Address]Universe)

	for _, node := range discovery.Nodes {
		for _, outputPort := range node.config.OutputPorts {
			universes[outputPort.Address] = Universe{
				controller: controller,
				Address:    outputPort.Address,
			}
		}
	}

	return universes
}

func (controller *Controller) Universe(address Address) Universe {
	return Universe{
		controller: controller,
		Address:    address,
	}
}

type Universe struct {
	controller *Controller

	Address Address
}

func (universe Universe) WriteDMX(dmx dmx.Universe) error {
	return universe.controller.SendDMX(universe.Address, dmx)
}

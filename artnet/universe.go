package artnet

import (
	"fmt"
	"io"
	"time"

	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/logging"
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

// XXX: unsafe
func (controller *Controller) Universe(address Address) *Universe {
	var universe = controller.universes[address]

	if universe == nil {
		universe = &Universe{
			log:        controller.config.Log.Logger("universe", address),
			controller: controller,
			address:    address,

			ticker:    time.NewTicker(controller.config.DMXRefresh),
			writeChan: make(chan dmx.Universe),
		}

		go universe.run()

		controller.universes[address] = universe
	}

	return universe
}

type Universe struct {
	log        logging.Logger
	controller *Controller
	address    Address

	ticker    *time.Ticker
	writeChan chan dmx.Universe

	dmx dmx.Universe
}

func (universe *Universe) Address() Address {
	return universe.address
}

func (universe *Universe) String() string {
	return fmt.Sprintf("artnet=%v", universe.address)
}

// any further WriteDMX() will fail
// XXX: delete from controller?
func (universe *Universe) close() {
	var writeChan chan dmx.Universe

	writeChan, universe.writeChan = universe.writeChan, nil

	close(writeChan)
}

func (universe *Universe) run() {
	defer universe.ticker.Stop()

	for {
		select {
		case dmx, ok := <-universe.writeChan:
			if !ok {
				break
			}

			universe.dmx = dmx

			universe.log.Debugf("Update: len=%d", len(universe.dmx))

		case <-universe.ticker.C:
			universe.log.Debugf("Refresh: len=%d", len(universe.dmx))
		}

		if err := universe.controller.SendDMX(universe.address, universe.dmx); err != nil {
			universe.log.Warnf("SendDMX: %v", err)
		}
	}
}

// Update DMX output
// Triggers immediate refresh + further refresh sends
//
// Returns io.EOF if closed
func (universe *Universe) WriteDMX(dmx dmx.Universe) error {
	if universe.writeChan == nil {
		return io.EOF
	}

	universe.writeChan <- dmx.Copy()

	return nil
}

// Stop the refresh cycle
//
// any further WriteDMX() will fail
func (universe *Universe) Close() {
	universe.close()
}

package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
)

type Output struct {
	log *log.Entry
	dmx dmx.Universe

	dmxWriter dmx.Writer // or nil

	dirty bool
}

func (output *Output) init(dmxWriter dmx.Writer) {
	output.dmxWriter = dmxWriter
}

func (output *Output) Get(address dmx.Address) dmx.Channel {
	return output.dmx.Get(address)
}

func (output *Output) Set(address dmx.Address, value dmx.Channel) {
	output.dirty = true
	output.dmx.Set(address, value)
}

func (output *Output) refresh() error {
	if output.dmxWriter == nil {
		return nil
	}

	output.log.Debugf("Output len=%v writer=%v:", len(output.dmx), output.dmxWriter)
	output.log.Debug(output.dmx)

	if err := output.dmxWriter.WriteDMX(output.dmx); err != nil {
		return err
	}

	output.dirty = false

	return nil
}

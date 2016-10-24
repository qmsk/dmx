package heads

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
)

type Output struct {
	log *log.Entry

	universe Universe
	dmx      dmx.Universe

	dmxWriter dmx.Writer // or nil

	dirty bool
}

func (output *Output) String() string {
	return fmt.Sprintf("%d", output.universe)
}

func (output *Output) init(dmxWriter dmx.Writer) {
	output.dmxWriter = dmxWriter
}

func (output *Output) GetDMX(address dmx.Address) dmx.Channel {
	return output.dmx.Get(address)
}

func (output *Output) GetValue(address dmx.Address) Value {
	return Value(output.GetDMX(address)) / 255.0
}

func (output *Output) SetDMX(address dmx.Address, value dmx.Channel) {
	output.dirty = true
	output.dmx.Set(address, value)
}

// Set value 0.0 .. 1.0
//
// Returns value at DMX percision
func (output *Output) SetValue(address dmx.Address, value Value) Value {
	var dmxChannel = dmx.Channel(value * 255.0)

	output.SetDMX(address, dmxChannel)

	return Value(dmxChannel) / 255.0
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

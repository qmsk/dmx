package heads

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
	"github.com/SpComb/qmsk-dmx/artnet"
	"github.com/qmsk/go-web"
)

type OutputConfig struct {
	Universe   Universe
	ArtNetNode *artnet.NodeConfig
}

type outputMap map[Universe]*Output

type Output struct {
	log *log.Entry

	config   OutputConfig
	universe Universe
	dmx      dmx.Universe

	dmxWriter dmx.Writer // or nil

	dirty bool
}

func (output *Output) String() string {
	return fmt.Sprintf("%d", output.universe)
}

func (output *Output) init(config OutputConfig, dmxWriter dmx.Writer) {
	output.config = config
	output.universe = config.Universe
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

func (outputMap outputMap) makeAPI() []APIOutput {
	var apiOutputs []APIOutput

	for _, output := range outputMap {
		apiOutputs = append(apiOutputs, output.makeAPI())
	}

	return apiOutputs
}

func (outputMap outputMap) GetREST() (web.Resource, error) {
	return outputMap.makeAPI(), nil
}

// Web API
type APIOutput struct {
	OutputConfig
}

func (output *Output) makeAPI() APIOutput {
	return APIOutput{
		OutputConfig: output.config,
	}
}

func (output *Output) GetREST() (web.Resource, error) {
	return output.makeAPI(), nil
}

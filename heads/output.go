package heads

import (
	"fmt"
	"time"

	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

type OutputConfig struct {
	Seen    time.Time
	Address string
	Port    int

	Artnet interface{} // metadata
}

type outputMap map[Universe]*Output

func (outputMap outputMap) makeAPI() APIOutputs {
	var apiOutputs = make(APIOutputs)

	for _, output := range outputMap {
		apiOutputs[output.String()] = output.makeAPI()
	}

	return apiOutputs
}

func (outputMap outputMap) GetREST() (web.Resource, error) {
	return outputMap.makeAPI(), nil
}

type OutputConnection struct {
	time      time.Time
	config    OutputConfig
	dmxWriter dmx.Writer // or nil
}

type Output struct {
	log    logging.Logger
	events Events

	universe   Universe
	dmx        dmx.Universe
	connection *OutputConnection
}

func (output *Output) String() string {
	return fmt.Sprintf("%d", output.universe)
}

func (output *Output) init(logger logging.Logger, universe Universe) {
	output.log = logger.Logger("universe", universe)
	output.universe = universe
	output.dmx = dmx.MakeUniverse()
}

func (output *Output) connect(config OutputConfig, dmxWriter dmx.Writer) {
	if output.connection == nil {
		output.connection = &OutputConnection{
			time:      time.Now(),
			config:    config,
			dmxWriter: dmxWriter,
		}
	} else {
		output.connection.config = config
		output.connection.dmxWriter = dmxWriter
	}

	output.apply()
}

func (output *Output) apply() {
	var id = output.String()

	output.events.update(APIEvents{
		Outputs: APIOutputs{
			id: output.makeAPI(),
		},
	})
}

func (output *Output) GetDMX(address dmx.Address) dmx.Channel {
	return output.dmx.Get(address)
}

func (output *Output) GetValue(address dmx.Address) Value {
	return Value(output.GetDMX(address)) / 255.0
}

func (output *Output) SetDMX(address dmx.Address, value dmx.Channel) {
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

func (output *Output) Refresh() error {
	// XXX: not goroutine-safe vs connect()
	if output.connection == nil {
		return nil
	}

	output.log.Debugf("Output len=%v writer=%v:", len(output.dmx), output.connection.dmxWriter)
	output.log.Debug(output.dmx)

	if err := output.connection.dmxWriter.WriteDMX(output.dmx); err != nil {
		return err
	}

	return nil
}

// Web API
type APIOutputs map[string]APIOutput

type APIOutput struct {
	Universe  Universe
	Connected *time.Time

	*OutputConfig
}

func (output *Output) makeAPI() APIOutput {
	var apiOutput = APIOutput{
		Universe: output.universe,
	}

	if output.connection != nil {
		apiOutput.Connected = &output.connection.time
		apiOutput.OutputConfig = &output.connection.config
	}

	return apiOutput
}

func (output *Output) GetREST() (web.Resource, error) {
	return output.makeAPI(), nil
}

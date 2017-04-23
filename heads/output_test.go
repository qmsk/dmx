package heads

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/logging"
)

var testLogger = logging.New("heads-test")

type testOutputWriter struct {
	mock.Mock
}

func (testOutputWriter testOutputWriter) WriteDMX(dmxUniverse dmx.Universe) error {
	return testOutputWriter.Called(dmxUniverse).Error(0)
}

func TestGetSet(t *testing.T) {
	var output Output
	output.init(testLogger, Universe(0))

	output.SetDMX(dmx.Address(1), dmx.Channel(0))
	assert.Equal(t, output.GetDMX(dmx.Address(1)), dmx.Channel(0))
	assert.Equal(t, output.GetValue(dmx.Address(1)), Value(0.0))

	output.SetDMX(dmx.Address(1), dmx.Channel(255))
	assert.Equal(t, output.GetDMX(dmx.Address(1)), dmx.Channel(255))
	assert.Equal(t, output.GetValue(dmx.Address(1)), Value(1.0))

	output.SetValue(dmx.Address(1), Value(0.0))
	assert.Equal(t, output.GetDMX(dmx.Address(1)), dmx.Channel(0))
	assert.Equal(t, output.GetValue(dmx.Address(1)), Value(0.0))

	output.SetValue(dmx.Address(1), Value(1.0))
	assert.Equal(t, output.GetDMX(dmx.Address(1)), dmx.Channel(255))
	assert.Equal(t, output.GetValue(dmx.Address(1)), Value(1.0))
}

func TestAPIOutputsNotConnected(t *testing.T) {
	var output Output
	var outputs = outputMap{}

	output.init(testLogger, Universe(1))
	outputs[output.universe] = &output

	assert.Equal(t, outputs.makeAPI(), APIOutputs{
		"1": APIOutput{
			Universe:  1,
			Connected: nil,
		},
	})
}

func TestAPIOutputNotConnected(t *testing.T) {
	var output Output
	output.init(testLogger, Universe(1))

	assert.Equal(t, output.makeAPI(), APIOutput{
		Universe:  1,
		Connected: nil,
	})
}

func TestRefreshNotConnected(t *testing.T) {
	var events = new(testEvents)

	var output = Output{
		events: events,
	}
	output.init(testLogger, Universe(1))

	assert.NoError(t, output.Refresh())
}

func TestAPIOutputConnected(t *testing.T) {
	var seen = time.Now()
	var writer testOutputWriter
	var events = new(testEvents)

	var config = OutputConfig{
		Seen:    seen,
		Address: "test",
		Port:    1,
	}
	var output = Output{
		events: events,
	}
	output.init(testLogger, Universe(1))

	events.On("update").Return()

	output.connect(config, writer)

	var api = output.makeAPI()
	var connected = *api.Connected

	assert.IsType(t, connected, time.Now())
	assert.Equal(t, api, APIOutput{
		Universe:     1,
		Connected:    &connected,
		OutputConfig: &config,
	})
}

func TestRefreshConnected(t *testing.T) {
	var seen = time.Now()
	var writer = new(testOutputWriter)
	var events = new(testEvents)

	var config = OutputConfig{
		Seen:    seen,
		Address: "test",
		Port:    1,
	}
	var output = Output{
		events: events,
	}
	output.init(testLogger, Universe(1))

	events.On("update").Return()

	output.connect(config, writer)

	writer.On("WriteDMX", dmx.Universe{}).Return(nil)

	assert.NoError(t, output.Refresh())
}

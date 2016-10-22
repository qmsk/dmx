package heads

import (
	log "github.com/Sirupsen/logrus"
	dmx "github.com/SpComb/qmsk-dmx"
)

type Universe int

type Output struct {
	log      *log.Entry
	universe Universe
	dmx      dmx.Universe
	writer   dmx.Writer

	dirty bool
}

func (output *Output) Get(address dmx.Address) dmx.Channel {
	return output.dmx.Get(address)
}

func (output *Output) Set(address dmx.Address, value dmx.Channel) {
	output.dirty = true
	output.dmx.Set(address, value)
}

func (output *Output) refresh() error {
	output.log.Debugf("Output len=%v writer=%v:", len(output.dmx), output.writer)
	output.log.Debug(output.dmx)

	if err := output.writer.WriteDMX(output.dmx); err != nil {
		return err
	}

	output.dirty = false

	return nil
}

type Heads struct {
	log     *log.Entry
	outputs map[Universe]*Output
	heads   []*Head
}

func New() *Heads {
	return &Heads{
		log:     log.WithField("package", "heads"),
		outputs: make(map[Universe]*Output),
	}
}

func (heads *Heads) Output(universe Universe, dmxWriter dmx.Writer) {
	heads.outputs[universe] = &Output{
		log:    heads.log.WithField("universe", universe),
		dmx:    dmx.MakeUniverse(),
		writer: dmxWriter,
	}
}

func (heads *Heads) Head(universe Universe, address dmx.Address, config HeadConfig) *Head {
	var output = heads.outputs[universe]
	var head = Head{
		config:  config,
		address: address,
	}

	head.init(output, config)

	heads.heads = append(heads.heads, &head)

	return &head
}

func (heads *Heads) Refresh() error {
	var refreshErr error

	for _, output := range heads.outputs {
		if err := output.refresh(); err != nil {
			refreshErr = err
		}
	}

	return refreshErr
}

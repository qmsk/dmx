package heads

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-dmx"
)

type Options struct {
	LibraryPath string `long:"heads-library" value-name:"PATH"`
}

func (options Options) Heads(config *Config) (*Heads, error) {
	var heads = Heads{
		log:     log.WithField("package", "heads"),
		outputs: make(map[Universe]*Output),
	}

	for _, headConfig := range config.Heads {
		if headType, exists := config.headTypes[headConfig.Type]; !exists {
			return nil, fmt.Errorf("Invalid Head.Type=%v", headConfig.Type)
		} else {
			headConfig.headType = headType
		}

		heads.addHead(headConfig)
	}

	return &heads, nil
}

type Heads struct {
	log     *log.Entry
	outputs map[Universe]*Output
	heads   []*Head
}

func (heads *Heads) output(universe Universe) *Output {
	output := heads.outputs[universe]
	if output == nil {
		output = &Output{
			log: heads.log.WithField("universe", universe),
			dmx: dmx.MakeUniverse(),
		}

		heads.outputs[universe] = output
	}

	return output
}

// Patch output
func (heads *Heads) Output(universe Universe, dmxWriter dmx.Writer) {
	heads.output(universe).init(dmxWriter)
}

// Patch head
func (heads *Heads) addHead(config HeadConfig) *Head {
	var output = heads.output(config.Universe)
	var head = Head{
		headType: config.headType,
		address:  config.Address,
	}

	head.init(output, config.headType)

	heads.heads = append(heads.heads, &head)

	return &head
}

func (heads *Heads) Each(fn func(head *Head)) {
	for _, head := range heads.heads {
		fn(head)
	}
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

package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-web"
)

func (heads *Heads) WebAPI() web.API {
	return web.MakeAPI(heads)
}

type API struct {
	Outputs []APIOutput        `json:"outputs"`
	Heads   map[string]APIHead `json:"heads"`
}

func (heads *Heads) Index(name string) (web.Resource, error) {
	log.Debugln("heads:Heads.Index", name)

	switch name {
	case "":
		return heads, nil
	case "outputs":
		return heads.outputs, nil
	case "heads":
		return heads.heads, nil
	default:
		return nil, nil
	}
}

func (heads *Heads) GetREST() (web.Resource, error) {
	log.Debug("heads:Heads.GetREST")
	return API{
		Outputs: heads.outputs.makeAPI(),
		Heads:   heads.heads.makeAPI(),
	}, nil
}

func (heads *Heads) Apply() error {
	log.Debug("heads:Heads.Apply")

	if err := heads.Refresh(); err != nil {
		log.Warnf("heads:Heads.Refresh: %v", err)

		return err
	}

	return nil
}

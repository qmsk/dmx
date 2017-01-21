package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qmsk/go-web"
)

type API struct {
	Outputs []APIOutput
	Heads   APIHeads
	Groups  APIGroups
	Presets presetMap
}

func (heads *Heads) WebAPI() web.API {
	return web.MakeAPI(heads)
}

func (heads *Heads) Index(name string) (web.Resource, error) {
	log.Debugln("heads:Heads.Index", name)

	switch name {
	case "":
		return heads, nil
	case "groups":
		return heads.groups, nil
	case "outputs":
		return heads.outputs, nil
	case "heads":
		return heads.heads, nil
	case "presets":
		return heads.presets, nil
	default:
		return nil, nil
	}
}

func (heads *Heads) GetREST() (web.Resource, error) {
	log.Debug("heads:Heads.GetREST")
	return API{
		Outputs: heads.outputs.makeAPI(),
		Heads:   heads.heads.makeAPI(),
		Groups:  heads.groups.makeAPI(),
		Presets: heads.presets,
	}, nil
}

func (heads *Heads) Apply() error {
	log.Debug("heads:Heads.Apply")

	// Refresh DMX output
	if err := heads.Refresh(); err != nil {
		log.Warnf("heads:Heads.Refresh: %v", err)

		return err
	}

	return nil
}

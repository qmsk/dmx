package heads

import (
	"net/http"

	"github.com/qmsk/go-web"
)

type API struct {
	Outputs APIOutputs
	Heads   APIHeads
	Groups  APIGroups
	Presets presetMap
}

func (heads *Heads) WebAPI() web.API {
	return web.MakeAPI(heads)
}

func (heads *Heads) Index(name string) (web.Resource, error) {
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

func (heads *Heads) makeAPI() API {
	return API{
		Outputs: heads.outputs.makeAPI(),
		Heads:   heads.heads.makeAPI(),
		Groups:  heads.groups.makeAPI(),
		Presets: heads.presets,
	}
}

func (heads *Heads) GetREST() (web.Resource, error) {
	return heads.makeAPI(), nil
}

func (heads *Heads) Apply() error {
	heads.log.Info("Apply")

	// Refresh DMX output
	if err := heads.Refresh(); err != nil {
		heads.log.Warn("Refresh: ", err)

		return err
	}

	return nil
}

func (heads *Heads) WebConfigPreset() http.Handler {
	return httpConfigPreset{heads}
}

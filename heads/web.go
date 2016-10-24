package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/SpComb/qmsk-web"
)

func (heads *Heads) WebAPI() web.API {
	return web.MakeAPI(heads)
}

type API struct {
	Heads map[string]APIHead `json:"heads"`
}

func (heads *Heads) Index(name string) (web.Resource, error) {
	log.Debug("heads:Heads.Index", name)

	switch name {
	case "":
		return heads, nil
	case "heads":
		return heads.heads, nil
	default:
		return nil, nil
	}
}

func (heads *Heads) GetREST() (web.Resource, error) {
	log.Debug("heads:Heads.GetREST")
	return API{
		Heads: heads.heads.makeAPI(),
	}, nil
}

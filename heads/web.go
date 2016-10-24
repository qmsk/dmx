package heads

import (
	"github.com/qmsk/e2/web"
)

type API struct {
	Heads map[string]APIHead `json:"heads"`
}

func (heads *Heads) WebAPI() web.API {
	return web.MakeAPI(heads)
}

func (heads *Heads) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return heads, nil
	case "heads":
		return heads.heads, nil
	default:
		return nil, nil
	}
}

func (heads *Heads) Get() (interface{}, error) {
	return API{
		Heads: heads.heads.apiDict(),
	}, nil
}

type headList headMap

func (heads headList) Get() (interface{}, error) {
	var apiHeads []APIHead

	for _, head := range heads {
		apiHeads = append(apiHeads, head.makeAPI())
	}

	return apiHeads, nil
}

func (heads headMap) apiDict() map[string]APIHead {
	var apiHeads = make(map[string]APIHead)

	for headID, head := range heads {
		apiHeads[headID] = head.makeAPI()
	}
	return apiHeads
}

func (headMap headMap) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return headList(headMap), nil
	default:
		return headMap[name], nil
	}
}

func (headMap headMap) Get() (interface{}, error) {
	return headMap.apiDict(), nil
}

package heads

import (
	"github.com/qmsk/go-web"
)

type groupMap map[GroupID]*Group

type APIGroups map[GroupID]APIGroup

func (groupMap groupMap) makeAPI() APIGroups {
	apiGroups := make(APIGroups)

	for groupID, group := range groupMap {
		apiGroups[groupID] = group.makeAPI()
	}

	return apiGroups
}

func (groupMap groupMap) makeAPIList() (apiGroups []APIGroup) {
	for _, group := range groupMap {
		apiGroups = append(apiGroups, group.makeAPI())
	}

	return
}

func (groupMap groupMap) GetREST() (web.Resource, error) {
	return groupMap.makeAPI(), nil
}

func (groupMap groupMap) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return groupMap.makeAPIList(), nil
	default:
		return groupMap[GroupID(name)], nil
	}
}

type Group struct {
	id     GroupID
	config GroupConfig
	heads  headMap

	color *GroupColor
}

// make empty group
func makeGroup(id GroupID, config GroupConfig) *Group {
	return &Group{
		id:     id,
		config: config,
		heads:  make(headMap),
	}
}

func (group *Group) addHead(head *Head) {
	group.heads[head.id] = head
}

// initialize group parameters from heads
func (group *Group) init() {
	if groupColor := group.makeColor(); groupColor.exists() {
		group.color = &groupColor
	}
}

func (group *Group) makeColor() GroupColor {
	var groupColor = GroupColor{
		headColors: make(map[HeadID]HeadColor),
	}

	for headID, head := range group.heads {
		if headColor := head.parameters.Color; headColor != nil {
			groupColor.headColors[headID] = *headColor
		}
	}

	return groupColor
}

type APIGroup struct {
	GroupConfig
	ID    GroupID
	Heads []HeadID

	Color *APIColor
}

func (group *Group) makeAPIHeads() (heads []HeadID) {
	for headID, _ := range group.heads {
		heads = append(heads, headID)
	}
	return
}

func (group *Group) makeAPI() APIGroup {
	return APIGroup{
		GroupConfig: group.config,
		ID:          group.id,
		Heads:       group.makeAPIHeads(),
		Color:       group.color.makeAPI(),
	}
}

func (group *Group) GetREST() (web.Resource, error) {
	return group.makeAPI(), nil
}

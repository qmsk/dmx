package heads

import (
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

type groups map[api.GroupID]*Group

func (groups groups) makeAPI() api.Groups {
	apiGroups := make(api.Groups)

	for groupID, group := range groups {
		apiGroups[groupID] = group.makeAPI()
	}

	return apiGroups
}

func (groups groups) makeAPIList() (apiGroups []api.Group) {
	for _, group := range groups {
		apiGroups = append(apiGroups, group.makeAPI())
	}

	return
}

func (groups groups) GetREST() (web.Resource, error) {
	return groups.makeAPI(), nil
}

func (groups groups) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return groups.makeAPIList(), nil
	default:
		return groups[api.GroupID(name)], nil
	}
}

// Group
type Group struct {
	log    logging.Logger
	id     api.GroupID
	config api.GroupConfig
	heads  heads
	events Events

	intensity *GroupIntensity
	color     *GroupColor
	colors    api.Colors
}

func (group *Group) addHead(head *Head) {
	group.heads[head.id] = head
}

// initialize group parameters from heads
func (group *Group) init() {
	// reverse-mappings for apply updates
	for _, head := range group.heads {
		head.initGroup(group)
	}

	if groupIntensity := group.makeIntensity(); groupIntensity.exists() {
		group.intensity = &groupIntensity
	}

	if groupColor := group.makeColor(); groupColor.exists() {
		group.color = &groupColor
	}

	// merge head ColorMaps
	for _, head := range group.heads {
		if colorMap := head.headType.Colors; colorMap != nil {
			group.colors.Merge(colorMap)
		}
	}
}

func (group *Group) makeIntensity() GroupIntensity {
	var groupIntensity = GroupIntensity{
		heads: make(map[api.HeadID]HeadIntensity),
	}

	for headID, head := range group.heads {
		if headIntensity := head.parameters.Intensity; headIntensity != nil {
			groupIntensity.heads[headID] = *headIntensity
		}
	}

	return groupIntensity
}

func (group *Group) makeColor() GroupColor {
	var groupColor = GroupColor{
		headColors: make(map[api.HeadID]HeadColor),
	}

	for headID, head := range group.heads {
		if headColor := head.parameters.Color; headColor != nil {
			groupColor.headColors[headID] = *headColor
		}
	}

	return groupColor
}

func (group *Group) makeAPIHeads() []api.HeadID {
	var heads = make([]api.HeadID, 0)

	for headID, _ := range group.heads {
		heads = append(heads, headID)
	}
	return heads
}

func (group *Group) makeAPIParams() api.GroupParams {
	var params api.GroupParams

	if group.intensity != nil {
		var intensity = group.intensity.makeAPI()

		params.Intensity = &intensity
	}

	if group.color != nil {
		var color = group.color.makeAPI()

		params.Color = &color
	}

	return params
}

func (group *Group) makeAPI() api.Group {
	return api.Group{
		GroupConfig: group.config,
		ID:          group.id,
		Heads:       group.makeAPIHeads(),
		Colors:      group.colors,
		GroupParams: group.makeAPIParams(),
	}
}

func (group *Group) GetREST() (web.Resource, error) {
	return group.makeAPI(), nil
}
func (group *Group) PostREST() (web.Resource, error) {
	return &APIGroupParams{group: group}, nil
}

// Web API
type APIGroupParams struct {
	group     *Group
	Intensity *APIIntensity `json:",omitempty"`
	Color     *APIColor     `json:",omitempty"`
}

func (apiGroupParams APIGroupParams) Apply() error {
	if apiGroupParams.Intensity != nil {
		if err := apiGroupParams.Intensity.initGroup(apiGroupParams.group.intensity); err != nil {
			return web.RequestError(err)
		} else if err := apiGroupParams.Intensity.Apply(); err != nil {
			return err
		}
	}

	if apiGroupParams.Color != nil {
		if err := apiGroupParams.Color.initGroup(apiGroupParams.group.color); err != nil {
			return web.RequestError(err)
		} else if err := apiGroupParams.Color.Apply(); err != nil {
			return err
		}
	}

	return nil
}

// Web API Events
func (group *Group) Apply() error {
	group.log.Info("Apply")

	group.events.update(api.Event{
		Heads: group.heads.makeAPI(),
		Groups: api.Groups{
			group.id: group.makeAPI(),
		},
	})

	return nil
}

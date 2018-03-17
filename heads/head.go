package heads

import (
	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/api"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

type heads map[api.HeadID]*Head

func (heads heads) makeAPI() api.Heads {
	var apiHeads = make(api.Heads)

	for headID, head := range heads {
		apiHeads[headID] = head.makeAPI()
	}
	return apiHeads
}

func (heads heads) makeAPIList() []api.Head {
	var apiHeads = make([]api.Head, 0, len(heads))

	for _, head := range heads {
		apiHeads = append(apiHeads, head.makeAPI())
	}

	return apiHeads
}

func (heads heads) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return headsView{heads}, nil
	default:
		return heads[api.HeadID(name)], nil
	}
}

func (heads heads) GetREST() (web.Resource, error) {
	return heads.makeAPI(), nil
}

// GET /heads/ => []api.Head
type headsView struct {
	heads heads
}

func (view headsView) GetREST() (web.Resource, error) {
	return view.heads.makeAPIList(), nil
}

type HeadParameters struct {
	Intensity *HeadIntensity `json:"intensity,omitempty"`
	Color     *HeadColor     `json:"color,omitempty"`
}

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	log logging.Logger

	id       api.HeadID
	config   api.HeadConfig
	headType api.HeadType
	output   *Output
	events   Events
	groups   groups

	channels   channels
	parameters HeadParameters
}

func (head *Head) String() string {
	return string(head.id)
}

func (head *Head) Name() string {
	if head.config.Name != "" {
		return head.config.Name
	} else {
		return string(head.id)
	}
}

func (head *Head) init() {
	head.channels = make(channels)

	for index, channelConfig := range head.headType.Channels {
		var channel = &Channel{
			id:      channelConfig.ID(),
			config:  channelConfig,
			index:   uint(index),
			output:  head.output,
			address: dmx.Address(head.config.Address) + dmx.Address(index),
		}

		channel.init()

		head.channels[channel.id] = channel
	}

	// set parameters
	if headIntensity := head.Intensity(); headIntensity.exists() {
		head.parameters.Intensity = &headIntensity
	}
	if headColor := head.Color(); headColor.exists() {
		head.parameters.Color = &headColor
	}
}

// Head is member of Group
func (head *Head) initGroup(group *Group) {
	head.groups[group.id] = group
}

func (head *Head) Channel(config api.ChannelConfig) *Channel {
	return head.channels[config.ID()]
}

func (head *Head) Intensity() HeadIntensity {
	return HeadIntensity{
		channel: head.Channel(api.ChannelConfig{Intensity: true}),
	}
}

func (head *Head) Color() HeadColor {
	return HeadColor{
		red:       head.Channel(api.ChannelConfig{Color: api.ChannelColorRed}),
		green:     head.Channel(api.ChannelConfig{Color: api.ChannelColorGreen}),
		blue:      head.Channel(api.ChannelConfig{Color: api.ChannelColorBlue}),
		intensity: head.Channel(api.ChannelConfig{Intensity: true}),
	}
}

func (head *Head) Parameters() HeadParameters {
	return head.parameters
}

func (head *Head) makeAPI() api.Head {
	var apiHead = api.Head{
		ID:     head.id,
		Config: head.config,
		Type:   head.headType,

		Channels: head.channels.makeAPI(),
	}

	if head.parameters.Intensity != nil {
		var intensity = head.parameters.Intensity.makeAPI()

		apiHead.Intensity = &intensity
	}

	if head.parameters.Color != nil {
		var color = head.parameters.Color.makeAPI()

		apiHead.Color = &color
	}

	return apiHead
}

func (head *Head) GetREST() (web.Resource, error) {
	return head.makeAPI(), nil
}

// Web API POST
type APIHeadParams struct {
	head *Head

	Channels  map[string]api.ChannelParams `json:",omitempty"`
	Intensity *api.Intensity               `json:",omitempty"`
	Color     *api.Color                   `json:",omitempty"`
}

func (head *Head) PostREST() (web.Resource, error) {
	// parameters only, not configuration
	return &APIHeadParams{head: head}, nil
}

func (post *APIHeadParams) Apply() error {
	post.head.log.Info("Apply parameters: %#v", post)

	for channelID, channelParams := range post.Channels {
		if channel := post.head.channels.GetID(channelID); channel == nil {
			return web.Errorf(404, "Channel not found: %v", channelID)
		} else {
			channelParams.channel = channel
		}

		if err := channelParams.Apply(); err != nil {
			return err
		}
	}

	if post.Intensity != nil {
		if err := post.Intensity.initHead(post.head.parameters.Intensity); err != nil {
			return web.RequestError(err)
		} else if err := post.Intensity.Apply(); err != nil {
			return err
		}
	}

	if post.Color != nil {
		if err := post.Color.initHead(post.head.parameters.Color); err != nil {
			return web.RequestError(err)
		} else if err := post.Color.Apply(); err != nil {
			return err
		}
	}

	return nil
}

func (head *Head) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return head, nil
	case "channels":
		return head.channels, nil
	case "intensity":
		return head.parameters.Intensity, nil
	case "color":
		return head.parameters.Color, nil
	default:
		return nil, nil
	}
}

// Web API Events
func (head *Head) Apply() error {
	head.log.Info("Apply")

	head.events.update(APIEvents{
		Heads: APIHeads{
			head.id: head.makeAPI(),
		},
		Groups: head.groups.makeAPI(),
	})

	return nil
}

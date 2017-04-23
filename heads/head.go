package heads

import (
	"fmt"

	"github.com/qmsk/dmx"
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

// Config type
type TypeID string

type HeadType struct {
	Vendor string
	Model  string
	Mode   string
	URL    string

	Channels []ChannelType
	Colors   ColorMap
}

func (headType HeadType) String() string {
	return fmt.Sprintf("%v/%v=%v", headType.Vendor, headType.Model, headType.Mode)
}

func (headType HeadType) IsColor() bool {
	for _, channelType := range headType.Channels {
		if channelType.Color != "" {
			return true
		}
	}
	return false
}

// Config
type HeadID string

func (headID HeadID) index(index uint) HeadID {
	return HeadID(fmt.Sprintf("%s.%d", headID, index+1))
}

type HeadConfig struct {
	Type     TypeID
	Universe Universe
	Address  dmx.Address
	Name     string
	Count    uint // Clone multiple copies of the head at id.N
	Groups   []GroupID

	headType *HeadType
}

// Number of channels used by head for count indexing
func (headConfig HeadConfig) step() uint {
	return uint(len(headConfig.headType.Channels))
}

// Return an indexed copy of the head, step addresses ahead
func (headConfig HeadConfig) index(index uint) HeadConfig {
	// copy
	var indexed HeadConfig = headConfig

	indexed.Address = indexed.Address + dmx.Address(index*headConfig.step())

	return indexed
}

// Top-level map
type headMap map[HeadID]*Head

type APIHeads map[HeadID]APIHead

func (heads headMap) makeAPI() APIHeads {
	var apiHeads = make(APIHeads)

	for headID, head := range heads {
		apiHeads[headID] = head.makeAPI()
	}
	return apiHeads
}

type headList headMap

func (heads headList) GetREST() (web.Resource, error) {
	var apiHeads []APIHead

	for _, head := range heads {
		apiHeads = append(apiHeads, head.makeAPI())
	}

	return apiHeads, nil
}

func (headMap headMap) Index(name string) (web.Resource, error) {
	switch name {
	case "":
		return headList(headMap), nil
	default:
		return headMap[HeadID(name)], nil
	}
}

func (headMap headMap) GetREST() (web.Resource, error) {
	return headMap.makeAPI(), nil
}

// Channels
type HeadChannels map[ChannelType]*Channel

func (headChannels HeadChannels) GetID(id string) *Channel {
	for channelType, channel := range headChannels {
		if channelType.String() == id {
			return channel
		}
	}

	return nil
}

func (headChannels HeadChannels) makeAPI() APIChannels {
	var apiChannels = make(APIChannels)

	for channelType, channel := range headChannels {
		apiChannels[channelType.String()] = channel.makeAPI()
	}

	return apiChannels
}

func (headChannels HeadChannels) GetREST() (web.Resource, error) {
	return headChannels.makeAPI(), nil
}

func (headChannels HeadChannels) Index(name string) (web.Resource, error) {
	if channel := headChannels.GetID(name); channel == nil {
		return nil, nil
	} else {
		return web.GetPostResource(channel), nil
	}
}

type HeadParameters struct {
	Intensity *HeadIntensity `json:"intensity,omitempty"`
	Color     *HeadColor     `json:"color,omitempty"`
}

// A single DMX receiver using multiple consecutive DMX channels from a base address within a single universe
type Head struct {
	log logging.Logger

	id       HeadID
	config   HeadConfig
	headType *HeadType
	output   *Output
	events   Events
	groups   groupMap

	channels   HeadChannels
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
	head.channels = make(HeadChannels)

	for channelIndex, channelType := range head.headType.Channels {
		var channel = &Channel{
			channelType: channelType,
			index:       uint(channelIndex),
			output:      head.output,
			address:     head.config.Address + dmx.Address(channelIndex),
		}

		channel.init()

		head.channels[channelType] = channel
	}

	// set parameters
	if headIntensity := head.getIntensity(); headIntensity.exists() {
		head.parameters.Intensity = &headIntensity
	}
	if headColor := head.getColor(); headColor.exists() {
		head.parameters.Color = &headColor
	}
}

// Head is member of Group
func (head *Head) initGroup(group *Group) {
	head.groups[group.id] = group
}

func (head *Head) getChannel(channelType ChannelType) *Channel {
	return head.channels[channelType]
}

func (head *Head) getIntensity() HeadIntensity {
	return HeadIntensity{
		channel: head.getChannel(ChannelType{Intensity: true}),
	}
}

func (head *Head) getColor() HeadColor {
	return HeadColor{
		red:       head.getChannel(ChannelType{Color: ColorChannelRed}),
		green:     head.getChannel(ChannelType{Color: ColorChannelGreen}),
		blue:      head.getChannel(ChannelType{Color: ColorChannelBlue}),
		intensity: head.getChannel(ChannelType{Intensity: true}),
	}
}

func (head *Head) Parameters() HeadParameters {
	return head.parameters
}

// Web API GET
type APIHead struct {
	ID     HeadID
	Config HeadConfig
	Type   *HeadType

	Channels  map[string]APIChannel `json:",omitempty"`
	Intensity *APIIntensity         `json:",omitempty"`
	Color     *APIColor             `json:",omitempty"`
}

func (head *Head) makeAPI() APIHead {
	return APIHead{
		ID:     head.id,
		Config: head.config,
		Type:   head.headType,

		Channels:  head.channels.makeAPI(),
		Intensity: head.parameters.Intensity.makeAPI(),
		Color:     head.parameters.Color.makeAPI(),
	}
}

func (head *Head) GetREST() (web.Resource, error) {
	return head.makeAPI(), nil
}

// Web API POST
type APIHeadParams struct {
	head *Head

	Channels  map[string]APIChannelParams `json:",omitempty"`
	Intensity *APIIntensity               `json:",omitempty"`
	Color     *APIColor                   `json:",omitempty"`
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

package heads

import (
	"github.com/qmsk/dmx/logging"
	"github.com/qmsk/go-web"
)

// Websocket
func (heads *Heads) WebEvents() web.Events {
	heads.events.eventChan = make(chan web.Event)

	return web.MakeEvents(heads.events.eventChan)
}

type APIEvents struct {
	Heads  APIHeads
	Groups APIGroups
}

func (event *APIEvents) addHead(head *Head) {
	if event.Heads == nil {
		event.Heads = APIHeads{head.id: head.makeAPI()}
	} else {
		event.Heads[head.id] = head.makeAPI()
	}
}

func (event *APIEvents) addHeads(heads headMap) {
	if event.Heads == nil {
		event.Heads = heads.makeAPI()
	} else {
		for headID, head := range heads {
			event.Heads[headID] = head.makeAPI()
		}
	}
}

func (events *APIEvents) addGroup(group *Group) {
	if events.Groups == nil {
		events.Groups = APIGroups{group.id: group.makeAPI()}
	} else {
		events.Groups[group.id] = group.makeAPI()
	}
}

func (event *APIEvents) addGroups(groups groupMap) {
	if event.Groups == nil {
		event.Groups = groups.makeAPI()
	} else {
		for groupID, group := range groups {
			event.Groups[groupID] = group.makeAPI()
		}
	}
}

type Events struct {
	log       logging.Logger
	eventChan chan web.Event
}

func (events *Events) update(event APIEvents) {
	if events.eventChan != nil {
		events.log.Infof("update")
		events.eventChan <- event
	}
}

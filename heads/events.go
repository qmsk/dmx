package heads

import (
	"github.com/qmsk/go-web"
)

// Websocket
func (heads *Heads) WebEvents() web.Events {
	heads.events = &Events{
		eventChan: make(chan web.Event),
	}

	return web.MakeEvents(heads.events.eventChan)
}

type APIEvents struct {
	Heads  map[HeadID]APIHead
	Groups map[GroupID]APIGroup
}

type Events struct {
	eventChan chan web.Event
}

func (events *Events) updateHead(id HeadID, apiHead APIHead) {
	if events != nil {
		events.eventChan <- APIEvents{Heads: APIHeads{id: apiHead}}
	}
}
func (events *Events) updateGroup(id GroupID, apiGroup APIGroup) {
	if events != nil {
		events.eventChan <- APIEvents{Groups: APIGroups{id: apiGroup}}
	}
}

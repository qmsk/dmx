package heads

import (
	"github.com/SpComb/qmsk-dmx/logging"
	"github.com/qmsk/go-web"
)

// Websocket
func (heads *Heads) WebEvents() web.Events {
	heads.events = &Events{
		log:       heads.options.Log.Logger("events", nil),
		eventChan: make(chan web.Event),
	}

	return web.MakeEvents(heads.events.eventChan)
}

type APIEvents struct {
	Heads  map[HeadID]APIHead
	Groups map[GroupID]APIGroup
}

type Events struct {
	log       logging.Logger
	eventChan chan web.Event
}

func (events *Events) updateHead(id HeadID, apiHead APIHead) {
	if events != nil {
		events.log.Infof("update head %v", id)
		events.eventChan <- APIEvents{Heads: APIHeads{id: apiHead}}
	}
}
func (events *Events) updateGroup(id GroupID, apiGroup APIGroup) {
	if events != nil {
		events.log.Infof("update group %v", id)
		events.eventChan <- APIEvents{Groups: APIGroups{id: apiGroup}}
	}
}

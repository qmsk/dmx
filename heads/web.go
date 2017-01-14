package heads

import (
	log "github.com/Sirupsen/logrus"
	"github.com/qmsk/go-web"
)

type API struct {
	Outputs []APIOutput
	Heads   map[HeadID]APIHead
	Groups  map[GroupID]APIGroup
}

func (heads *Heads) Index(name string) (web.Resource, error) {
	log.Debugln("heads:Heads.Index", name)

	switch name {
	case "":
		return heads, nil
	case "groups":
		return heads.groups, nil
	case "outputs":
		return heads.outputs, nil
	case "heads":
		return heads.heads, nil
	default:
		return nil, nil
	}
}

func (heads *Heads) GetREST() (web.Resource, error) {
	log.Debug("heads:Heads.GetREST")
	return API{
		Outputs: heads.outputs.makeAPI(),
		Heads:   heads.heads.makeAPI(),
		Groups:  heads.groups.makeAPI(),
	}, nil
}

func (heads *Heads) Apply() error {
	log.Debug("heads:Heads.Apply")

	// API Events
	heads.events.flush()

	// DMX output
	if err := heads.Refresh(); err != nil {
		log.Warnf("heads:Heads.Refresh: %v", err)

		return err
	}

	return nil
}

// Websocket
func (heads *Heads) WebEvents() chan web.Event {
	heads.events.eventChan = make(chan web.Event)

	return heads.events.eventChan
}

type APIEvents struct {
	Heads map[string]APIHead
}

type Events struct {
	eventChan chan web.Event
	pending   APIEvents
}

func (events *Events) init() {
	events.pending = APIEvents{
		Heads: make(map[string]APIHead),
	}
}

// Add pending event for Head
// XXX: unsafe
func (events *Events) updateHead(id string, apiHead APIHead) {
	events.pending.Heads[id] = apiHead
}

// Write out and clear pending events
// XXX: unsafe
func (events *Events) flush() {
	events.eventChan <- events.pending
	events.init()
}

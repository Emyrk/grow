package events

import (
	"github.com/emyrk/grow/world"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

type EventController struct {
	log zerolog.Logger

	newEvents      chan Event
	existingEvents map[uint64]Event

	eventOrder []uint64
	eventsEnd  int
}

func NewEventController(logger zerolog.Logger) *EventController {
	return &EventController{
		newEvents:      make(chan Event, 10000),
		existingEvents: make(map[uint64]Event),
		log:            logger,
	}
}

func (w *EventController) SendEvent(e Event) error {
	select {
	case w.newEvents <- e:
	default:
		return xerrors.Errorf("event queue full, rejected %d", e.ID())
	}
	return nil
}

func (ec *EventController) UpdateInOrder(w *world.World, clean bool) {
	var cleaned []uint64
	// Process all events in the given event order.
	// If an event id is 0, that means it is skipped.
	for idx, id := range ec.eventOrder {
		if id == 0 {
			continue
		}
		// Tick the event and get a new one
		c, err := ec.existingEvents[id].Tick(w)
		if err != nil {
			AddLogFields(ec.log.Error(), c).
				Err(err).
				Msg("update existing event")
		}
		// If it is nil, we delete the event from the ones we are tracking
		if c == nil {
			delete(ec.existingEvents, id)
			ec.eventOrder[idx] = 0
		} else {
			ec.existingEvents[id] = c
			if clean {
				cleaned = append(cleaned, c.ID())
			}
		}
	}

NewEventLoop:
	for {
		select {
		case e := <-ec.newEvents:
			c, err := e.Tick(w)
			if err != nil {
				AddLogFields(ec.log.Error(), c).
					Err(err).
					Msg("update new event")
			}
			if c != nil {
				ec.existingEvents[c.ID()] = e
				if clean {
					cleaned = append(cleaned, c.ID())
				} else {
					ec.eventOrder = append(ec.eventOrder, c.ID())
				}
			}
		default:
			break NewEventLoop
		}
	}

	// Removed all excess
	if clean {
		ec.eventOrder = cleaned
	}
}

func (ec *EventController) Update(w *world.World) {
	// Always clean until it's a problem
	ec.UpdateInOrder(w, true)

}

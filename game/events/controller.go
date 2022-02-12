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
		log:            logger.With().Str("service", "game").Logger(),
	}
}

// SendEvent will send the event to be queued to be played.
func (w *EventController) SendEvent(e Event) error {
	if e == nil {
		return nil
	}
	select {
	case w.newEvents <- e:
	default:
		return xerrors.Errorf("event queue full, rejected %d", e.GetID())
	}
	return nil
}

func (ec *EventController) UpdateInOrder(w *world.World, gametick uint64) (bool, []Event) {
	syncTick := SyncTick(gametick)
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
			if syncTick {
				cleaned = append(cleaned, c.GetID())
			}
		}
	}

	var newEvents []Event
	if syncTick {
	NewEventLoop:
		for {
			select {
			case e := <-ec.newEvents:
				newEvents = append(newEvents, e)
				c, err := e.Tick(w)
				if err != nil {
					AddLogFields(ec.log.Error(), c).
						Err(err).
						Msg("update new event")
				}
				if c != nil {
					ec.existingEvents[c.GetID()] = e
					cleaned = append(cleaned, c.GetID())
				}
			default:
				break NewEventLoop
			}
		}

	}

	// Removed all excess
	if syncTick {
		ec.eventOrder = cleaned
	}
	return syncTick, newEvents
}

func (ec *EventController) Update(w *world.World, gametick uint64) (bool, []Event) {
	sync, events := ec.UpdateInOrder(w, gametick)
	return sync, events
}

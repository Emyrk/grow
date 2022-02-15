package events

import (
	world2 "github.com/emyrk/grow/game/world"
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

func (ec *EventController) ReplaceEventList(evts []Event) {
	// Drain all new events
NewEventDrain:
	for {
		select {
		case <-ec.newEvents:
		default:
			break NewEventDrain
		}
	}

	// Reset the lists and maps
	ec.existingEvents = make(map[uint64]Event)
	ec.eventOrder = make([]uint64, 0, len(evts))
	for i := range evts {
		evt := evts[i]
		ec.eventOrder = append(ec.eventOrder, evt.GetID())
		ec.existingEvents[evt.GetID()] = evt
	}
}

func (ec *EventController) EventList() []Event {
	list := make([]Event, 0, len(ec.existingEvents))
	for _, v := range ec.eventOrder {
		if v == 0 {
			continue
		}
		list = append(list, ec.existingEvents[v])
	}
	return list
}

func (ec *EventController) UpdateInOrder(w *world2.World, gametick uint64) (bool, []Event) {
	syncTick := SyncTick(gametick)
	var cleaned []uint64
	// Process all events in the given event order.
	// If an event id is 0, that means it is skipped.
	for idx, id := range ec.eventOrder {
		if id == 0 {
			continue
		}
		// Tick the event and get a new one
		c, err := ec.existingEvents[id].Tick(gametick, w)
		if err != nil {
			AddLogFields(ec.log.Error(), c).
				Err(err).
				Msg("update existing event")
		}
		// If it is nil, we delete the event from the ones we are tracking
		if c == nil {
			//ec.log.Info().Uint64("eid", id).Msg("delete event")
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
				if _, ok := ec.existingEvents[e.GetID()]; ok {
					ec.log.Warn().Uint64("eid", e.GetID()).Msg("duplicate event")
					continue
				}
				newEvents = append(newEvents, e)
				c, err := e.Tick(gametick, w)
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

func (ec *EventController) Update(w *world2.World, gametick uint64) (bool, []Event) {
	sync, events := ec.UpdateInOrder(w, gametick)
	return sync, events
}

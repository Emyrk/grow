package events

import (
	"github.com/emyrk/grow/world"
	"github.com/rs/zerolog"
)

type EventType string

const (
	LeftClickEvent = "left-click"
)

type Event interface {
	ID() uint64
	Type() EventType
	// Tick will allow the event to advance 1 tick. If the event is done, it should return a nil.
	Tick(w *world.World) (Event, error)
}

func AddLogFields(l *zerolog.Event, e Event) *zerolog.Event {
	return l.Uint64("id", e.ID()).Str("type", string(e.Type())).Str("log_type", "game_event")
}

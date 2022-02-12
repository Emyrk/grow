package events

import (
	"github.com/emyrk/grow/internal/crand"
)

type baseEvent struct {
	ID uint64
}

func newBaseEvent() baseEvent {
	return baseEvent{
		ID: crand.Uint64(),
	}
}

func (e baseEvent) GetID() uint64 {
	return e.ID
}

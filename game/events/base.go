package events

import "math/rand"

type baseEvent struct {
	ID uint64
}

func newBaseEvent() baseEvent {
	return baseEvent{
		ID: rand.Uint64(),
	}
}

func (e baseEvent) GetID() uint64 {
	return e.ID
}

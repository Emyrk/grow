package events

import "math/rand"

type baseEvent struct {
	id uint64
}

func newBaseEvent() baseEvent {
	return baseEvent{
		id: rand.Uint64(),
	}
}

func (e baseEvent) ID() uint64 {
	return e.id
}

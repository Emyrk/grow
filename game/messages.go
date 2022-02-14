package game

import (
	"github.com/emyrk/grow/game/events"
	world2 "github.com/emyrk/grow/game/world"
)

type GameMessageType = string

const (
	// MsgTickEventList is the list of events that were accepted
	MsgTickEventList GameMessageType = "tick-event-list"
	// MsgGameSync is the full game sync if the payload exists. If the payload is empty, it is a client
	// asking for a sync
	MsgGameSync GameMessageType = "game-sync"
	// MsgGameNewEvents are new events from a client
	MsgGameNewEvents GameMessageType = "new-events"
)

// TickEventList is the list of new events for a given game tick.
type TickEventList struct {
	GameTick  uint64
	Eventlist events.EventList
}

func (TickEventList) Type() GameMessageType {
	return MsgTickEventList
}

// GameSync is all the data you need to instantly go to a game tick
type GameSync struct {
	World *world2.World
	// Present list of game events at the game tick in the correct order to be processed for THIS game tick
	EventList events.EventList
	GameTick  uint64
}

func (GameSync) Type() GameMessageType {
	return MsgGameSync
}

// NewEvents is the list of new events for a given game tick.
type NewEvents struct {
	Player    world2.PlayerID
	Eventlist events.EventList
}

func (NewEvents) Type() GameMessageType {
	return MsgGameNewEvents
}

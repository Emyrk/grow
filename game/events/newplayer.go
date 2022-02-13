package events

import (
	"image/color"

	"github.com/emyrk/grow/world"
)

type PlayerJoin struct {
	baseEvent
	PlayerID world.PlayerID
	Color    color.RGBA
	Team     uint16
}

func (c *PlayerJoin) Type() EventType {
	return PlayerJoined
}

func (c *PlayerJoin) Tick(w *world.World) (Event, error) {
	p := &world.Player{
		ID:    c.PlayerID,
		Color: c.Color,
		Team:  0,
	}
	w.Players.AddPlayer(p)
	return nil, nil
}

package events

import (
	world2 "github.com/emyrk/grow/game/world"
	"image/color"
)

type PlayerJoin struct {
	baseEvent
	PlayerID world2.PlayerID
	Color    color.RGBA
	Team     uint16
}

func (c *PlayerJoin) Type() EventType {
	return PlayerJoined
}

func (c *PlayerJoin) Tick(w *world2.World) (Event, error) {
	p := &world2.Player{
		ID:    c.PlayerID,
		Color: c.Color,
		Team:  0,
	}
	w.Players.AddPlayer(p)
	return nil, nil
}

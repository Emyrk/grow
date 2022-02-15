package events

import (
	"image/color"

	world2 "github.com/emyrk/grow/game/world"
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

func (c *PlayerJoin) Tick(_ uint64, w *world2.World) (Event, error) {
	p := &world2.Player{
		ID:    c.PlayerID,
		Color: c.Color,
		Team:  0,
	}
	w.Players.AddPlayer(p)
	x := p.ID % uint64(w.MapWidth)
	y := p.ID % uint64(w.MapHeight)
	w.PlayerStart(p.ID, int(x), int(y))
	return nil, nil
}

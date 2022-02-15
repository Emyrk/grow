package events

import (
	"image"

	world2 "github.com/emyrk/grow/game/world"
)

type ClickEvent struct {
	baseEvent
	Player   *world2.Player
	Pos      image.Point
	Current  int
	Duration int
	Skip     int
	Skipped  int

	Left bool
}

func NewClickEvent(player *world2.Player, left bool, x, y int) *ClickEvent {
	return &ClickEvent{
		baseEvent: newBaseEvent(),
		Player:    player,
		Pos: image.Point{
			X: x,
			Y: y,
		},
		Duration: 5,
		Skip:     0,
		Left:     left,
	}
}

func (c *ClickEvent) Type() EventType {
	return MouseClickEvent
}

func (c *ClickEvent) Tick(gametick uint64, w *world2.World) (Event, error) {
	if _, ok := w.Players[c.Player.ID]; !ok {
		return nil, nil
	}

	if c.Skipped > 0 {
		c.Skipped--
		return c, nil
	}

	c.Current++
	c.Skipped = c.Skip
	if c.Current > c.Duration {
		return nil, nil
	}

	if c.Left {
		for i := 0; i < 200; i++ {
			w.Attack(c.Player.ID, c.Pos.X, c.Pos.Y, 1)
		}
		return c, nil
	}
	w.Attack(c.Player.ID, c.Pos.X, c.Pos.Y, 0)
	return nil, nil

	if c.Skipped > 0 {
		c.Skipped--
		return c, nil
	}

	l := c.Current * 2
	half := l / 2

	tlx := c.Pos.X - half
	tly := c.Pos.Y + half
	brx := c.Pos.X + half
	bry := c.Pos.Y - half

	for cx := tlx; cx < brx; cx++ {
		for cy := tly; cy > bry; cy-- {
			w.Claim(c.Player.ID, cx, cy)
		}
	}

	c.Current++
	c.Skipped = c.Skip
	if c.Current > c.Duration {
		return nil, nil
	}
	return c, nil
}

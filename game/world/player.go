package world

import (
	"image/color"

	"github.com/emyrk/grow/internal/crand"
)

type PlayerSet map[PlayerID]*Player

func NewPlayerSet() PlayerSet {
	return make(map[PlayerID]*Player)
}

func (s *PlayerSet) AddRandomPlayer() *Player {
	p := RandomPlayer()
	return s.AddPlayer(p)
}

func (s *PlayerSet) NewPlayer(team uint16, col color.RGBA) *Player {
	p := NewPlayer(team, col)
	return s.AddPlayer(p)
}

func (s *PlayerSet) AddPlayer(p *Player) *Player {
	for {
		if _, ok := (*s)[p.ID]; ok {
			p.ID = crand.Uint64()
			continue
		}
		break
	}
	(*s)[p.ID] = p
	return p
}

type PlayerID = uint64

type Player struct {
	ID    PlayerID
	Color color.RGBA
	// Team 0 is FFA
	Team uint16
}

func NewPlayer(team uint16, col color.RGBA) *Player {
	return &Player{
		ID:    crand.Uint64(),
		Team:  team,
		Color: col,
	}
}

func RandomPlayer() *Player {
	return NewPlayer(0, color.RGBA{
		R: crand.Uint8(),
		G: crand.Uint8(),
		B: crand.Uint8(),
		A: crand.Uint8(),
	})
}

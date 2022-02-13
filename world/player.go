package world

import (
	"image/color"
	"math/rand"
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
			p.ID = PlayerID(rand.Uint32())
			continue
		}
		break
	}
	(*s)[p.ID] = p
	return p
}

type PlayerID = uint16

type Player struct {
	ID    PlayerID
	Color color.RGBA
	// Team 0 is FFA
	Team uint16
}

func NewPlayer(team uint16, col color.RGBA) *Player {
	return &Player{
		ID:    PlayerID(rand.Uint32()),
		Team:  team,
		Color: col,
	}
}

func RandomPlayer() *Player {
	return NewPlayer(0, color.RGBA{
		R: uint8(rand.Uint32()),
		G: uint8(rand.Uint32()),
		B: uint8(rand.Uint32()),
		A: uint8(rand.Uint32()),
	})
}

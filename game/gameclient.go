package game

import (
	"github.com/rs/zerolog"
)

// GameClient is able to run a game, but relies on a server to handle event ordering.
type GameClient struct {
	G        *Game
	Gametick uint64

	Log zerolog.Logger
}

func NewGameClient(log zerolog.Logger, cfg GameConfig) *GameClient {
	c := &GameClient{
		G:   NewGame(log, cfg),
		Log: log,
	}
	return c
}

func (c *GameClient) Update() error {
	// TODO: Handle blocking for events
	c.G.Update(c.Gametick)
	c.Gametick++
	return nil
}

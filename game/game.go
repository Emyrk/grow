package game

import (
	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/world"
	"github.com/rs/zerolog"
)

type ProcessEvents func(gametick uint64, events []events.Event)

type GameConfig struct {
	Players world.PlayerSet
	Width   int
	Height  int
}

type Game struct {
	World *world.World
	EC    *events.EventController
	Log   zerolog.Logger
}

func NewGame(log zerolog.Logger, cfg GameConfig) *Game {
	return &Game{
		World: world.NewWorld(cfg.Width, cfg.Height, cfg.Players),
		EC:    events.NewEventController(log),
		Log:   log.With().Str("service", "game").Logger(),
	}
}

// Update is called every 1/60 of a second
func (g *Game) Update(gametick uint64) (bool, []events.Event) {
	processEvents, evts := g.EC.Update(g.World, gametick)
	g.World.Update()
	return processEvents, evts
}

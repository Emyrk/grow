package game

import (
	"time"

	"github.com/emyrk/grow/game/events"
	world2 "github.com/emyrk/grow/game/world"
	"github.com/rs/zerolog"
)

type BroadcastGameMessage func(msgType GameMessageType, data []byte)

type GameConfig struct {
	Players world2.PlayerSet
	Width   int
	Height  int
}

type Game struct {
	// Started indicates the game started. No one else can join
	Started bool

	World *world2.World
	EC    *events.EventController
	Log   zerolog.Logger

	expSec time.Time
}

func NewGame(log zerolog.Logger, cfg GameConfig) *Game {
	return &Game{
		World: world2.NewWorld(cfg.Width, cfg.Height, cfg.Players),
		EC:    events.NewEventController(log),
		Log:   log.With().Str("service", "game").Logger(),
	}
}

// Update is called every 1/60 of a second
func (g *Game) Update(gametick uint64) (bool, []events.Event) {
	if gametick%60 == 0 {
		g.Log.Info().Float64("dur", time.Since(g.expSec).Seconds()).Msg("exp 1 sec")
		g.expSec = time.Now()
	}
	processEvents, evts := g.EC.Update(g.World, gametick)
	g.World.Update()
	return processEvents, evts
}

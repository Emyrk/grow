package game

import (
	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/game/keybinds"
	"github.com/emyrk/grow/world"
	"github.com/rs/zerolog"
)

type Game struct {
	World    *world.World
	ec       *events.EventController
	keybinds *keybinds.KeyWatcher
	pixels   []byte
	log      zerolog.Logger
}

func NewGame(log zerolog.Logger, width, height int, players world.PlayerSet, me *world.Player) *Game {
	return &Game{
		World:    world.NewWorld(width, height, players),
		ec:       events.NewEventController(log),
		keybinds: keybinds.NewKeybinds(me),
	}
}

func (g *Game) Update() error {
	actions := g.keybinds.Update()
	for i := range actions {
		err := g.ec.SendEvent(actions[i])
		if err != nil {
			events.AddLogFields(g.log.Error(), actions[i]).
				Err(err).
				Msg("send event")
		}
	}
	g.ec.Update(g.World)
	g.World.Update()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.World.Width(), g.World.Height()
}

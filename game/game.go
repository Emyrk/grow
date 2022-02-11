package game

import (
	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/game/keybinds"
	"github.com/emyrk/grow/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rs/zerolog"
)

type Game struct {
	world    *world.World
	ec       *events.EventController
	keybinds *keybinds.KeyWatcher
	pixels   []byte
	log      zerolog.Logger
}

func NewGame(log zerolog.Logger, width, height int, players world.PlayerSet, me *world.Player) *Game {
	return &Game{
		world:    world.NewWorld(width, height, players),
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
	g.ec.Update(g.world)
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, g.world.Width()*g.world.Height()*4)
	}
	g.world.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.world.Width(), g.world.Height()
}

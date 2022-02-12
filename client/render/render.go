package render

import (
	"context"

	"github.com/emyrk/grow/client/keybinds"
	"github.com/emyrk/grow/game"
	"github.com/emyrk/grow/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameRender struct {
	*game.GameClient

	keyWatcher *keybinds.KeyWatcher
	pixels     []byte
}

func NewGameRenderer(g *game.GameClient, me *world.Player) *GameRender {
	return &GameRender{
		GameClient: g,
		keyWatcher: keybinds.NewKeybinds(me),
	}
}

func (g *GameRender) Update() error {
	// Watch for new user generated events.
	actions := g.keyWatcher.Update()
	if len(actions) > 0 {
		err := g.SendGameEvents(context.Background(), actions)
		if err != nil {
			g.Log.
				Err(err).
				Int("event_count", len(actions)).
				Msg("send evts")
		}
	}
	err := g.GameClient.Update()
	if err != nil {
		return err
	}
	return nil
}

func (g *GameRender) Draw(screen *ebiten.Image) {
	gme := g.G
	if g.pixels == nil {
		g.pixels = make([]byte, gme.World.Width()*gme.World.Height()*4)
	}
	gme.World.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)
}

func (g *GameRender) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.G.World.Width(), g.G.World.Height()
}

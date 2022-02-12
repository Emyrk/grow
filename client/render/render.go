package render

import (
	"github.com/emyrk/grow/client/keybinds"
	"github.com/emyrk/grow/game"
	"github.com/emyrk/grow/game/events"
	"github.com/emyrk/grow/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameRender struct {
	game.GameClient

	keyWatcher *keybinds.KeyWatcher
	pixels     []byte
}

func NewGameRenderer(g game.GameClient, me *world.Player) *GameRender {
	return &GameRender{
		GameClient: g,
		keyWatcher: keybinds.NewKeybinds(me),
	}
}

func (g *GameRender) Update() error {
	// Watch for new user generated events.
	// TODO: @emyrk these should push to the server, not the game directly
	actions := g.keyWatcher.Update()
	for i := range actions {
		err := g.G.EC.SendEvent(actions[i])
		if err != nil {
			events.AddLogFields(g.Log.Error(), actions[i]).
				Err(err).
				Msg("send event")
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

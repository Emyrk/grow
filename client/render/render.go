package render

import (
	"github.com/emyrk/grow/game"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameRender struct {
	game.Game

	pixels []byte
}

func NewGameRenderer(g game.Game) *GameRender {
	return &GameRender{
		Game: g,
	}
}

func (g *GameRender) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, g.World.Width()*g.World.Height()*4)
	}
	g.World.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)
}

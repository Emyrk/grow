package testdata

import (
	"image/color"

	world2 "github.com/emyrk/grow/game/world"

	"github.com/emyrk/grow/game"
)

const (
	screenWidth  = 2048
	screenHeight = 2048
)

type TestGameData struct {
	Me      *world2.Player
	GameCfg game.GameConfig
}

func TestGame() TestGameData {
	me := world2.NewPlayer(0, color.RGBA{
		// 844a93
		R: 0x84,
		G: 0x4a,
		B: 0x93,
		A: 0xff,
	})
	players := world2.NewPlayerSet()
	me = players.AddPlayer(me)

	return TestGameData{
		Me: me,
		GameCfg: game.GameConfig{
			Players: players,
			Width:   screenWidth,
			Height:  screenHeight,
		},
	}
}

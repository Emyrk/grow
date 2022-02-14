package world

import (
	"image/color"

	"github.com/rs/zerolog"
)

type World struct {
	// MapTiles never change after initialized. It is the underlying map
	MapTiles []tile

	// PlayerTiles are who owns a tile
	PlayerTiles []PlayerID

	Players   PlayerSet
	MapWidth  int
	MapHeight int

	logger zerolog.Logger
}

func NewWorld(width, height int, players PlayerSet) *World {
	w := &World{
		PlayerTiles: make([]PlayerID, width*height),
		MapTiles:    make([]tile, width*height),
		MapWidth:    width,
		MapHeight:   height,
		Players:     players,
	}

	return w
}

// Claim sets the tile to the player's GetID
func (w *World) Claim(x, y int, playerID PlayerID) {
	if x < 0 || y < 0 || x > w.MapWidth || y > w.MapHeight {
		return
	}

	w.PlayerTiles[w.PointI(x, y)] = playerID
}

func (w *World) Update() {
	// Random stuff
	//x, y := rand.Intn(w.MapWidth), rand.Intn(w.MapHeight)
	//w.PlayerTiles[w.PointI(x, y)] = TileLand
}

func (w *World) PointXY(i int) (int, int) {
	x := i % w.MapWidth
	y := i / w.MapWidth
	return x, y
}

func (w *World) PointI(x, y int) int {
	// Find the height first
	i := y * w.MapWidth
	i += x
	return i
}

func (w *World) Width() int {
	return w.MapWidth
}

func (w *World) Height() int {
	return w.MapHeight
}

func (w *World) Draw(pix []byte) {
	for i, v := range w.PlayerTiles {
		var tileColor = color.RGBA{}
		if v == 0 {
			switch w.MapTiles[i] {
			case TileBlocked:
				tileColor = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
			case TileWater:
				tileColor = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
			case TileLand:
				tileColor = color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
			}
		} else {
			player, ok := w.Players[v]
			if ok {
				tileColor = player.Color
			} else {
				w.logger.Error().Msg("player tile with non-existent player")
			}
		}
		// R
		pix[4*i] = tileColor.R
		// G
		pix[4*i+1] = tileColor.G
		// B
		pix[4*i+2] = tileColor.B
		// A
		pix[4*i+3] = tileColor.A
	}
}

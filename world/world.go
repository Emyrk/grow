package world

import (
	"image/color"

	"github.com/rs/zerolog"
)

type World struct {
	area    []tile
	players PlayerSet
	width   int
	height  int
	logger  zerolog.Logger
}

// https://ebiten.org/examples/life.html
func NewWorld(width, height int, players PlayerSet) *World {
	w := &World{
		area:    make([]tile, width*height),
		width:   width,
		height:  height,
		players: players,
	}

	return w
}

// Claim sets the tile to the player's ID
func (w *World) Claim(x, y int, playerID PlayerID) {
	w.area[w.PointI(x, y)] = tile(playerID)
}

func (w *World) Update() {
	// Random stuff
	//x, y := rand.Intn(w.width), rand.Intn(w.height)
	//w.area[w.PointI(x, y)] = TileLand
}

func (w *World) PointXY(i int) (int, int) {
	x := i % w.width
	y := i / w.width
	return x, y
}

func (w *World) PointI(x, y int) int {
	// Find the height first
	i := y * w.width
	i += x
	return i
}

func (w *World) Width() int {
	return w.width
}

func (w *World) Height() int {
	return w.height
}

func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		var tileColor = color.RGBA{}
		switch v {
		case TileBlocked:
			tileColor = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		case TileWater:
			tileColor = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
		case TileLand:
			tileColor = color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
		default:
			pID := v
			if v < 0 {
				pID = v * -1
			}
			player := w.players[PlayerID(pID)]
			tileColor = player.Color
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

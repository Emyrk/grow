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

	Players PlayerSet
	// PlayerBorders is a lookup for the playerID. Their borders with other players
	PlayerBorders map[PlayerID]map[PlayerID][]int

	MapWidth  int
	MapHeight int

	logger zerolog.Logger
}

func NewWorld(width, height int, players PlayerSet) *World {
	w := &World{
		PlayerTiles:   make([]PlayerID, width*height),
		MapTiles:      make([]tile, width*height),
		PlayerBorders: make(map[PlayerID]map[PlayerID][]int),
		MapWidth:      width,
		MapHeight:     height,
		Players:       players,
	}

	for pid, _ := range players {
		w.PlayerBorders[pid] = make(map[PlayerID][]int)
	}

	return w
}

// PlayerStart starts a player
func (w *World) PlayerStart(playerID PlayerID, x, y int) {
	w.Claim(playerID, x, y)
}

func (w *World) Attack(att PlayerID, x, y int, strength int) {
	borderPxs := w.PlayerBorders[att]
	if borderPxs == nil {
		return // No attack
	}

	idx := w.PointI(x, y)
	if !w.SafeTile(idx) {
		return
	}

	def := w.PlayerTiles[w.PointI(x, y)]
	w.attack(att, def, strength)
}

func (w *World) attack(att PlayerID, def PlayerID, strength int) int {
	borderPxs := w.PlayerBorders[att]
	if borderPxs == nil {
		return 0 // No attack
	}

	defPx := w.PlayerBorders[att][def]
	if len(defPx) == 0 {
		return 0
	}

	var buckets [1000][]int
	for _, px := range defPx {
		if w.PlayerTiles[px] != def {
			continue // Don't add pixels that are wrong
		}
		var friends int
		neighbors := w.aroundXY(w.PointXY(px))
		for _, v := range neighbors {
			if !w.SafeTile(v) {
				continue
			}
			if w.PlayerTiles[v] == att {
				friends++
			}
		}
		w.MapTiles[px] = tile(friends * 2)
		buckets[friends] = append(buckets[friends], px)
	}

	if strength == 0 {
		return 0
	}

	// Refresh the borders after each attack for now
	w.PlayerBorders[att][def] = make([]int, 0)
	var claimed int
	var consumed int
	friends := len(buckets) - 1
	left := make([]int, 0)
	for ; friends >= 0; friends-- {
		consumed = 0
		next := buckets[friends]
		if claimed >= strength {
			left = append(left, next...)
			continue
		}
		for _, px := range next {
			consumed++
			fnd := w.PlayerTiles[px]
			if fnd != def {
				continue
			}

			ax, ay := w.PointXY(px)
			w.Claim(att, ax, ay)
			claimed++
			if claimed >= strength {
				left = make([]int, len(next)-consumed)
				copy(left, next[consumed:])
				break
			}
		}
	}

	w.PlayerBorders[att][def] = append(left, w.PlayerBorders[att][def]...)
	return w.attack(att, def, strength-claimed)
}

// Claim sets the tile to the player's GetID
func (w *World) Claim(playerID PlayerID, x, y int) {
	if x < 0 || y < 0 || x > w.MapWidth || y > w.MapHeight {
		return
	}

	idx := w.PointI(x, y)
	w.PlayerTiles[idx] = playerID
	w.addBorders(playerID, x, y)
}

func (w *World) addBorders(playerID PlayerID, x, y int) {
	borderPx := w.aroundXY(x, y)
	ourPx := w.PointI(x, y)
	for _, bIdx := range borderPx {
		if bIdx < 0 || bIdx >= len(w.PlayerTiles) {
			continue
		}
		bpid := w.PlayerTiles[bIdx]
		if bpid == playerID {
			continue
		}
		w.addBorder(playerID, bpid, bIdx)
		// We also now border them
		if bpid != 0 {
			w.addBorder(bpid, playerID, ourPx)
		}
	}
}

func (w *World) addBorder(player, border PlayerID, idx int) {
	if w.PlayerBorders[player] == nil {
		w.PlayerBorders[player] = make(map[PlayerID][]int)
	}
	if w.PlayerBorders[player][border] == nil {
		w.PlayerBorders[player][border] = make([]int, 0)
	}

	w.PlayerBorders[player][border] = append(w.PlayerBorders[player][border], idx)
}

//func (w *World) removeBorder(player, border PlayerID, idx int) {
//	if w.PlayerBorders[player] == nil {
//		w.PlayerBorders[player] = make(map[PlayerID][]int)
//	}
//	if w.PlayerBorders[player][border] == nil {
//		w.PlayerBorders[player][border] = make([]int, 0)
//	}
//	delete(w.PlayerBorders[player][border], idx)
//}

func (w *World) aroundXY(x, y int) []int {
	return []int{
		w.PointI(x, y-1), // North
		w.PointI(x+1, y), // East
		w.PointI(x, y+1), // South
		w.PointI(x-1, y), // West

		w.PointI(x+1, y-1), // NE
		w.PointI(x-1, y-1), // NW

		w.PointI(x+1, y+1), // SE
		w.PointI(x-1, y+1), // SW
	}
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

func (w *World) SafeTile(idx int) bool {
	if idx < 0 || idx >= len(w.PlayerTiles) {
		return false
	}
	return true
}

func (w *World) Draw(pix []byte) {
	for i, v := range w.PlayerTiles {
		var tileColor = color.RGBA{}
		if v == 0 {
			mv := w.MapTiles[i]
			switch mv {
			case TileBlocked:
				tileColor = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
			case TileWater:
				tileColor = color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
			case TileLand:
				tileColor = color.RGBA{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
			default:
				mu := uint8(mv) * 10
				tileColor = color.RGBA{R: mu, G: mu, B: mu, A: 0xff}
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

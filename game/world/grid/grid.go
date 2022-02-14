package grid

import "github.com/emyrk/grow/internal/crand"

type Grid struct {
	Width  int
	Height int

	Shapes []*Shape
}

func NewGrid(width, height int) *Grid {
	return &Grid{
		Width:  width,
		Height: height,
	}
}

func (g *Grid) AddShape(s *Shape) {
	s.ID = crand.Uint64()
	g.Shapes = append(g.Shapes, s)
}

func (g *Grid) PointXY(i int) (int, int) {
	x := i % g.Width
	y := i / g.Height
	return x, y
}

func (g *Grid) PointI(x, y int) int {
	i := y * g.Width
	i += x
	return i
}

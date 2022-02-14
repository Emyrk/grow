package grid

import (
	"image"
	"image/color"

	"github.com/emyrk/grow/internal/crand"
)

type Shape struct {
	ID           uint64
	BoundingRect image.Rectangle
	Color        color.RGBA

	globalPoints []image.Point
	localPoints  []image.Point

	// hullLocalPoints are needed for drawing
	hullLocalPoints []image.Point
}

// AddPoint takes a global coordinate
func (s *Shape) AddPoint(x, y int) {
	s.globalPoints = append(s.globalPoints, image.Point{X: x, Y: y})

	// Bounding
	if x < s.BoundingRect.Min.X {
		s.BoundingRect.Min.X = x
	} else if x > s.BoundingRect.Max.X {
		s.BoundingRect.Max.X = x
	}

	if y < s.BoundingRect.Min.Y {
		s.BoundingRect.Min.Y = y
	} else if x > s.BoundingRect.Max.Y {
		s.BoundingRect.Max.Y = y
	}

	s.updatePoints()
}

// RawPoints returns the points on the global grid
func (s *Shape) RawPoints() []image.Point {
	return s.globalPoints
}

// LocalPoints returns the points as they are in the local bounding rectangle
func (s *Shape) LocalPoints() []image.Point {
	return s.localPoints
}

func (s *Shape) updatePoints() {
	// Update local points
	s.localPoints = make([]image.Point, 0, len(s.globalPoints))
	for _, gpt := range s.globalPoints {
		pt := s.localPoint(gpt.X, gpt.Y)
		s.localPoints = append(s.localPoints, pt)
	}
}

func (s *Shape) localPoint(x, y int) image.Point {
	pt := image.Point{X: x, Y: y}
	pt.X -= s.BoundingRect.Min.X
	pt.Y -= s.BoundingRect.Min.Y
	return pt
}

func NewRectangle(r image.Rectangle) *Shape {
	s := new(Shape)
	defer s.updatePoints()
	s.globalPoints = []image.Point{
		image.Point{
			X: r.Min.X,
			Y: r.Max.Y,
		},
		image.Point{
			X: r.Max.X,
			Y: r.Max.Y,
		}, image.Point{
			X: r.Max.X,
			Y: r.Min.Y,
		}, image.Point{
			X: r.Min.X,
			Y: r.Min.Y,
		},
	}
	s.BoundingRect = r

	return s
}

func NewDiamond(center image.Point, size int) *Shape {
	base := size / 2
	s := new(Shape)
	defer s.updatePoints()
	s.globalPoints = []image.Point{
		// Top point
		{
			X: center.X,
			Y: center.Y + base,
		},
		{

			X: center.X + base,
			Y: center.Y,
		},
		{
			X: center.X,
			Y: center.Y - base,
		},
		{
			X: center.X - base,
			Y: center.Y,
		},
	}
	s.Color = color.RGBA{
		R: crand.Uint8(),
		G: crand.Uint8(),
		B: crand.Uint8(),
		A: 0xff,
	}
	s.BoundingRect = image.Rectangle{
		Min: image.Point{
			X: center.X - base,
			Y: center.Y - base,
		},
		Max: image.Point{
			X: center.X + base,
			Y: center.Y + base,
		},
	}

	return s
}

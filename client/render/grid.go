package render

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/rs/zerolog"

	"github.com/emyrk/grow/game/world/grid"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type shape struct {
	S        *grid.Shape
	DrawOpts *ebiten.DrawImageOptions
	Img      *ebiten.Image
}

// GridRender will render grids... yea that is it.
type GridRender struct {
	*grid.Grid
	log zerolog.Logger

	cached   map[uint64]*shape
	viewPort *ebiten.GeoM
}

func NewGridRenderer(log zerolog.Logger, g *grid.Grid) *GridRender {
	return &GridRender{
		Grid:     g,
		log:      log.With().Str("service", "grid_render").Logger(),
		cached:   make(map[uint64]*shape),
		viewPort: &ebiten.GeoM{},
	}
}

var last time.Time

func (g *GridRender) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.viewPort.Translate(-1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.viewPort.Translate(1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.viewPort.Translate(0, -1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.viewPort.Translate(0, 1)
	}
	dx, dy := ebiten.Wheel()
	if dx != 0 || dy != 0 {
		fmt.Println(dx, dy)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		sh := grid.NewShape([]image.Point{
			{X: x, Y: y},
		})

		//sh := grid.NewDiamond(image.Point{
		//	X: x,
		//	Y: y,
		//}, 50)
		g.AddShape(sh)
		g.log.Info().Msg("Draw diamond")
	}

	// Add a point to the last shape
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if len(g.Shapes)-1 >= 0 {
			s := g.Shapes[len(g.Shapes)-1]
			x, y := ebiten.CursorPosition()
			s.AddPoint(x, y)
			delete(g.cached, s.ID)
		}
		g.log.Info().Msg("Add point")
	}

	last = time.Now()
	return nil
}

func (g *GridRender) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

var i int

// Render is mainly to help debugging
func (g *GridRender) Draw(screen *ebiten.Image) {
	for _, s := range g.Shapes {
		sImg, ok := g.cached[s.ID]
		if !ok {
			dx, dy := s.BoundingRect.Dx(), s.BoundingRect.Dy()
			if dx <= 0 || dy <= 0 {
				continue
			}
			canvas := ebiten.NewImage(s.BoundingRect.Dx(), s.BoundingRect.Dy())
			gCtx := gg.NewContextForImage(canvas)

			gCtx.SetColor(s.Color)
			pts := s.LocalPoints()
			startPt := pts[0]
			gCtx.MoveTo(float64(startPt.X), float64(startPt.Y))
			for _, pt := range pts {
				gCtx.LineTo(float64(pt.X), float64(pt.Y))
			}
			gCtx.LineTo(float64(startPt.X), float64(startPt.Y))
			gCtx.SetFillStyle(gg.NewSolidPattern(s.Color))
			gCtx.Fill()

			gCtx.SetColor(color.White)
			for _, pt := range pts {
				gCtx.DrawPoint(float64(pt.X), float64(pt.Y), 5)
				gCtx.Fill()
			}
			img := ebiten.NewImageFromImage(gCtx.Image())

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(s.BoundingRect.Min.X), float64(s.BoundingRect.Min.Y))
			sImg = &shape{
				S:        s,
				DrawOpts: opts,
				Img:      img,
			}
			g.cached[s.ID] = sImg
		}
		opts := *g.viewPort
		opts.Translate(float64(s.BoundingRect.Min.X), float64(s.BoundingRect.Min.Y))
		screen.DrawImage(sImg.Img, &ebiten.DrawImageOptions{
			GeoM: opts,
		})
	}
}

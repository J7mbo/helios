package helios

import "time"

type Point struct {
	x float64
	y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func (p *Point) GetX() float64 {
	return p.x
}

func (p *Point) GetY() float64 {
	return p.y
}

type Region struct {
	topLeft *Point
	width   int
	height  int

	screen *Screen
	finder *Finder
}

func NewRegion(topLeft *Point, width, height int, screen *Screen, finder *Finder) *Region {
	return &Region{topLeft, width, height, screen, finder}
}

func (r *Region) GetTopLeft() *Point {
	return r.topLeft
}

func (r *Region) GetWidth() int {
	return r.width
}

func (r *Region) GetHeight() int {
	return r.height
}

func (r *Region) Highlight(t time.Duration) {
	r.screen.highlighter.Highlight(&HighlightRequest{
		ScreenWidth:  r.screen.width,
		ScreenHeight: r.screen.height,
		X:            r.topLeft.x,
		Y:            r.topLeft.y,
		Width:        float64(r.width),
		Height:       float64(r.height),
		Duration:     t.Seconds(),
	})
}

func (r *Region) Find(i *Image) *Match {
	return r.screen.finder.Find(i, r)
}

func (r *Region) FindAll(i *Image) []*Match {
	return r.screen.finder.FindAll(i, r)
}

func (r *Region) Wait(i *Image, t time.Duration) *Match {
	return r.finder.Wait(i, r, t)
}

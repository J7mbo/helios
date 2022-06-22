package helios

import (
	"github.com/go-vgo/robotgo"
	"time"
)

type Screen struct {
	Region

	highlighter *Highlighter
	finder      *Finder
	interval    *PollInterval
}

func NewScreen(highlighter *Highlighter, interval *PollInterval) *Screen {
	w, h := robotgo.GetScreenSize()

	screen := &Screen{}
	screen.Region.topLeft = &Point{0, 0}
	screen.Region.width = w
	screen.Region.height = h
	screen.highlighter = highlighter

	finder := NewFinder(screen, interval)
	screen.finder = finder

	return screen
}

func (s *Screen) NewRegion(topLeft *Point, width, height float64) *Region {
	return NewRegion(topLeft, int(width), int(height), s, s.finder)
}

func (s *Screen) Highlight(t time.Duration) {
	s.highlighter.Highlight(&HighlightRequest{
		ScreenWidth:  s.width,
		ScreenHeight: s.height,
		X:            s.topLeft.x,
		Y:            s.topLeft.y,
		Width:        float64(s.width),
		Height:       float64(s.height),
		Duration:     t.Seconds(),
	})
}

func (s *Screen) Find(i *Image) *Match {
	return s.finder.Find(i, &s.Region)
}

func (s *Screen) FindAll(i *Image) []*Match {
	return s.finder.FindAll(i, &s.Region)
}

func (s *Screen) Wait(i *Image, t time.Duration) *Match {
	return s.finder.Wait(i, &s.Region, t)
}

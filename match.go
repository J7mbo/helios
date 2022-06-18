package helios

import (
	"time"
)

type Match struct {
	screen      *Screen
	img         *Image
	highlighter *Highlighter
	clicker     *Clicker

	x          float64
	y          float64
	confidence float32
}

func NewMatch(img *Image, x float64, y float64, confidence float32, screen *Screen, highlighter *Highlighter) *Match {
	return &Match{
		screen:      screen,
		img:         img,
		x:           x,
		y:           y,
		highlighter: highlighter,
		confidence:  confidence,
	}
}

func (m *Match) GetConfidence() float32 {
	return m.confidence
}

func (m *Match) Highlight(t time.Duration) {
	if m == nil {
		return
	}

	bounds := m.img.img.Bounds()
	point := bounds.Size()

	m.highlighter.Highlight(&HighlightRequest{
		ScreenWidth:  m.screen.width,
		ScreenHeight: m.screen.height,
		X:            m.x,
		Y:            m.y,
		Width:        float64(point.X),
		Height:       float64(point.Y),
		Duration:     t.Seconds(),
	})
}

func (m *Match) Click() *Match {
	m.clicker.Click(m)

	return m
}

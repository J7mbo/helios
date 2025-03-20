package helios

import (
	"time"
)

type Match struct {
	Region

	screen      *Screen
	img         *Image
	highlighter *Highlighter
	clicker     *Clicker

	confidence float32
}

func NewMatch(img *Image, confidence float32, screen *Screen, highlighter *Highlighter, region *Region) *Match {
	return &Match{
		screen:      screen,
		img:         img,
		Region:      *region,
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

	m.highlighter.Highlight(&HighlightRequest{
		ScreenWidth:  m.screen.width,
		ScreenHeight: m.screen.height,
		X:            m.topLeft.x,
		Y:            m.topLeft.y,
		Width:        float64(m.width),
		Height:       float64(m.height),
		Duration:     t.Seconds(),
	})
}

func (m *Match) Click() *Match {
	m.clicker.Click(m)

	return m
}

func (m *Match) MoveMouse() *Match {
	m.clicker.MoveMouseInRegion(&m.Region)

	return m
}

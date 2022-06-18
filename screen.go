package helios

import (
	"bytes"
	"github.com/go-vgo/robotgo"
	"gocv.io/x/gocv"
	"image/png"
	"time"
)

type Screen struct {
	width  int
	height int

	highlighter *Highlighter
}

func NewScreen(highlighter *Highlighter) *Screen {
	w, h := robotgo.GetScreenSize()

	screen := &Screen{width: w, height: h}
	screen.highlighter = highlighter

	return screen
}

func (s *Screen) Highlight(t time.Duration) {
	s.highlighter.Highlight(&HighlightRequest{
		ScreenWidth:  s.width,
		ScreenHeight: s.height,
		X:            0,
		Y:            0,
		Width:        float64(s.width * 2),
		Height:       float64(s.height * 2),
		Duration:     t.Seconds(),
	})
}

func (s *Screen) Find(i *Image) *Match {
	// @todo abstract to "finder" for reuse in screen and match
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, i.img)
	bts := buf.Bytes()

	screen := robotgo.CaptureScreen(0, 0, s.width, s.height)
	backgroundImg := robotgo.ToBitmapBytes(screen)

	img, err := gocv.IMDecode(backgroundImg, gocv.IMReadAnyColor)
	if err != nil {
		panic(err)
	}

	template, err := gocv.IMDecode(bts, gocv.IMReadAnyColor)
	if err != nil {
		panic(err)
	}

	matResult := gocv.NewMat()
	mask := gocv.NewMat()

	defer func() { _ = template.Close() }()
	defer func() { _ = img.Close() }()
	defer func() { _ = matResult.Close() }()
	defer func() { _ = mask.Close() }()

	gocv.MatchTemplate(img, template, &matResult, gocv.TmCcoeffNormed, mask)

	_, maxConfidence, _, maxLoc := gocv.MinMaxLoc(matResult)

	if maxConfidence < float32(i.confidenceThreshold) {
		return nil
	}

	return NewMatch(i, float64(maxLoc.X), float64(maxLoc.Y), maxConfidence, s, s.highlighter)
}

func (s *Screen) FindAll(i *Image) []*Match {
	// @todo abstract to "finder" for reuse in screen and match
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, i.img)
	bts := buf.Bytes()

	screen := robotgo.CaptureScreen(0, 0, s.width, s.height)
	backgroundImg := robotgo.ToBitmapBytes(screen)

	img, err := gocv.IMDecode(backgroundImg, gocv.IMReadAnyColor)
	if err != nil {
		panic(err)
	}

	template, err := gocv.IMDecode(bts, gocv.IMReadAnyColor)
	if err != nil {
		panic(err)
	}

	matResult := gocv.NewMat()
	mask := gocv.NewMat()

	gocv.MatchTemplate(img, template, &matResult, gocv.TmCcoeffNormed, mask)
	defer func() { _ = mask.Close() }()
	defer func() { _ = matResult.Close() }()

	//_, maxConfidence, _, maxLoc := gocv.MinMaxLoc(matResult)

	//for {
	//	_, maxConfidence, _, maxLoc := gocv.MinMaxLoc(matResult)
	//}

	//imgFindW := template.Cols()
	//imgFindH := template.Rows()
	//newMat := gocv.NewMat()

	var matches []*Match
	for y := 0; y < matResult.Rows(); y++ {
	Loop:
		for x := 0; x < matResult.Cols(); x++ {
			if matResult.GetFloatAt(y, x) >= float32(i.confidenceThreshold) {
				// Make sure we haven't already put something + / - 5px in results already
				for _, m := range matches {
					// Within 5 above
					// E.g. X = 125 and Y = 100. We have a saved match at X: 120, Y: 95.
					// Ignore because it's within 5px above.
					// If the match X
					if inBetween(float64(x), m.x-5, m.x+5) && inBetween(float64(y), m.y-5, m.y+5) {
						continue Loop
					}
					//if float64(x) <= m.x+5 && float64(y) <= m.y+5 {
					//	continue Loop
					//}

					//if inBetween(float64(x), m.x-5, m.x+5) {
					//	continue Loop
					//}
					//// Within 5 below
					//if float64(x) >= m.x-5 && float64(y) >= m.y-5 {
					//	continue Loop
					//}
				}

				matches = append(matches, NewMatch(i, float64(x), float64(y), matResult.GetFloatAt(y, x), s, s.highlighter))
			}
		}
	}

	//
	//if maxConfidence < float32(i.confidenceThreshold) {
	//	return nil
	//}
	return matches

	//return NewMatch(i, float64(maxLoc.X), float64(maxLoc.Y), maxConfidence, s, s.highlighter)
}

func inBetween(i, min, max float64) bool {
	return (i >= min) && (i <= max)
}

package helios

import (
	"bytes"
	"github.com/go-vgo/robotgo"
	"gocv.io/x/gocv"
	"image/png"
	"time"
)

type Finder struct {
	screen       *Screen
	pollInterval *PollInterval
}

func NewFinder(screen *Screen, pollInterval *PollInterval) *Finder {
	return &Finder{screen: screen, pollInterval: pollInterval}
}

func (f *Finder) Find(i *Image, r *Region) *Match {
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, i.img)
	bts := buf.Bytes()

	screen := robotgo.CaptureScreen(int(r.topLeft.x), int(r.topLeft.y), r.width, r.height)
	scaleSize := robotgo.ScaleF() // On a low res screen, this is 1, but on MBPs with retina displays, it's 2
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

	defer func() {
		_ = template.Close()
		_ = img.Close()
		_ = matResult.Close()
		_ = mask.Close()
	}()

	gocv.MatchTemplate(img, template, &matResult, gocv.TmCcoeffNormed, mask)

	_, maxConfidence, _, maxLoc := gocv.MinMaxLoc(matResult)

	if maxConfidence < float32(i.confidenceThreshold) {
		return nil
	}

	region := &Region{
		topLeft: &Point{
			x: float64((maxLoc.X / int(scaleSize)) + int(r.topLeft.x)),
			y: float64((maxLoc.Y / int(scaleSize)) + int(r.topLeft.y)),
		},
		width:  i.img.Bounds().Size().X / int(scaleSize),
		height: i.img.Bounds().Size().Y / int(scaleSize),
		screen: f.screen,
	}

	return NewMatch(i, maxConfidence, f.screen, f.screen.highlighter, region)
}

func (f *Finder) FindAll(i *Image, r *Region) []*Match {
	buf := new(bytes.Buffer)
	_ = png.Encode(buf, i.img)
	bts := buf.Bytes()

	screen := robotgo.CaptureScreen(int(r.topLeft.x), int(r.topLeft.y), r.width, r.height)
	scaleSize := robotgo.ScaleF() // On a low res screen, this is 1, but on MBPs with retina displays, it's 2
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

	defer func() {
		_ = template.Close()
		_ = img.Close()
		_ = matResult.Close()
		_ = mask.Close()
	}()

	gocv.MatchTemplate(img, template, &matResult, gocv.TmCcoeffNormed, mask)

	var matches []*Match
	for y := 0; y < matResult.Rows(); y++ {
	Loop:
		for x := 0; x < matResult.Cols(); x++ {
			if matResult.GetFloatAt(y, x) >= float32(i.confidenceThreshold) {
				region := &Region{
					topLeft: &Point{
						x: float64((x / int(scaleSize)) + int(r.topLeft.x)),
						y: float64((y / int(scaleSize)) + int(r.topLeft.y)),
					},
					width:  i.img.Bounds().Size().X / int(scaleSize),
					height: i.img.Bounds().Size().Y / int(scaleSize),
					screen: f.screen,
				}

				// Make sure we haven't already put something + / - 5px in results already.
				// Ignore shifts in > / < 5px from each found match, otherwise we'll have many duplicates.
				for _, m := range matches {
					if inBetween(region.topLeft.x, m.topLeft.x-5, m.topLeft.x+5) &&
						inBetween(region.topLeft.y, m.topLeft.y-5, m.topLeft.y+5) {
						continue Loop
					}
				}

				matches = append(
					matches,
					NewMatch(i, matResult.GetFloatAt(y, x), f.screen, f.screen.highlighter, region),
				)
			}
		}
	}

	return matches
}

func (f *Finder) Wait(i *Image, r *Region, t time.Duration) *Match {
	// Default to 0.25 seconds if config not provided for PollInterval.
	pollInterval := f.pollInterval
	if f.pollInterval == nil {
		pollInterval = &PollInterval{100 * time.Millisecond}
	}

	for {
		if time.Now().Unix() > time.Now().Add(t).Unix() {
			break
		}

		if match := f.Find(i, r); match != nil {
			return match
		}

		time.Sleep(pollInterval.Duration)
	}

	return nil
}

func inBetween(i, min, max float64) bool {
	return (i >= min) && (i <= max)
}

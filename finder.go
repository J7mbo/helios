package helios

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
	"os"
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
	// Currently only works for the main monitor.
	backgroundImg := robotgo.CaptureImg()
	templateImagePath := "./template.png"

	// Because of bug: https://github.com/vcaesar/gcv/issues/3, we have to save the screenshot first
	// and then use it.
	robotgo.SavePng(backgroundImg, templateImagePath)

	templateImg, _, err := robotgo.DecodeImg(templateImagePath)
	if err != nil {
		// @todo return error too?
		fmt.Println(err)
		return nil
	}

	if err := os.Remove(templateImagePath); err != nil {
		fmt.Println(err)
		return nil
	}

	subImg, _, err := robotgo.DecodeImg(i.GetPath())
	if err != nil {
		return nil
	}

	// This just calls OpenCV MinMaxLoc()
	// These two give different results, we need both....
	_, maxConfidence, _, maxLoc := gcv.FindImg(subImg, templateImg)

	if maxConfidence < float32(i.confidenceThreshold) {
		return nil
	}

	scaleSize := robotgo.ScaleF()

	region := &Region{
		topLeft: &Point{
			x: float64(maxLoc.X) / scaleSize,
			y: float64(maxLoc.Y) / scaleSize,
		},
		width:  i.img.Bounds().Size().X / int(scaleSize),
		height: i.img.Bounds().Size().Y / int(scaleSize),
		screen: f.screen,
	}

	return NewMatch(i, maxConfidence, f.screen, f.screen.highlighter, region)
}

func (f *Finder) FindAll(i *Image, r *Region) []*Match {
	// Currently only works for the main monitor.
	backgroundImg := robotgo.CaptureImg()

	templateImagePath := "./template.png"

	// Because of bug: https://github.com/vcaesar/gcv/issues/3, we have to save the screenshot first
	// and then use it.
	robotgo.SavePng(backgroundImg, templateImagePath)

	templateImg, _, err := robotgo.DecodeImg(templateImagePath)
	if err != nil {
		// @todo return error too?
		fmt.Println(err)
		return nil
	}

	if err := os.Remove(templateImagePath); err != nil {
		fmt.Println(err)
		return nil
	}

	subImg, _, err := robotgo.DecodeImg(i.GetPath())
	if err != nil {
		return nil
	}

	results := gcv.FindAllImg(subImg, templateImg)
	scaleSize := robotgo.ScaleF()

	var matches []*Match

	for _, result := range results {
		maxConfidence := float32(result.MaxVal[0])
		if result.MaxVal[0] < maxConfidence {
			continue
		}

		region := &Region{
			topLeft: &Point{
				x: float64(result.TopLeft.X) / scaleSize,
				y: float64(result.TopLeft.Y) / scaleSize,
			},
			width:  i.img.Bounds().Size().X / int(scaleSize),
			height: i.img.Bounds().Size().Y / int(scaleSize),
			screen: f.screen,
		}

		matches = append(matches, NewMatch(i, maxConfidence, f.screen, f.screen.highlighter, region))
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

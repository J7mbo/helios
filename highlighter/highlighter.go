package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/akamensky/argparse"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/color"
	"os"
	"time"
)

type args struct {
	ScreenWidth  int
	ScreenHeight int
	X            float64
	Y            float64
	Width        float64
	Height       float64
	LineWidth    float64
	Duration     float64
	Colour       color.Color
}

func parseArgs() *args {
	parser := argparse.NewParser("highlighter", "Highlights areas of the screen for you")

	p := args{}

	screenWidth := parser.Int("", "screen-width", &argparse.Options{
		Required: true,
		Help:     "Full screen width in pixels",
	})

	screenHeight := parser.Int("", "screen-height", &argparse.Options{
		Required: true,
		Help:     "Full screen height in pixels",
	})

	x := parser.Float("", "x", &argparse.Options{
		Required: true,
		Help:     "X coordinate for the box to be drawn",
	})

	y := parser.Float("", "y", &argparse.Options{
		Required: true,
		Help:     "Y coordinate for the box to be drawn",
	})

	width := parser.Float("", "w", &argparse.Options{
		Required: true,
		Help:     "width of the box to be drawn",
	})

	height := parser.Float("", "h", &argparse.Options{
		Required: true,
		Help:     "height of the box to be drawn",
	})

	duration := parser.Float("", "d", &argparse.Options{
		Required: true,
		Help:     "duration to highlight for",
	})

	lineWidth := parser.Float("", "l", &argparse.Options{
		Required: false,
		Help:     "height of the box to be drawn",
		Default:  5.0,
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	p.ScreenWidth = *screenWidth
	p.ScreenHeight = *screenHeight
	p.X = *x
	p.Y = *y
	p.Height = *height
	p.Width = *width
	p.Duration = *duration
	p.LineWidth = *lineWidth

	// Default to red for now.
	p.Colour = color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}

	return &p
}

func main() {
	args := parseArgs()

	a := app.New()
	drv := a.Driver().(desktop.Driver)
	a.Settings().SetTheme(&transparentFyneTheme{})
	w := drv.CreateSplashWindow()

	w.Resize(fyne.Size{Width: float32(args.ScreenWidth), Height: float32(args.ScreenHeight)})

	img := canvas.NewImageFromImage(drawRectangle(args))
	img.FillMode = canvas.ImageFillStretch

	w.SetContent(img)

	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)

	go func() {
		time.Sleep(time.Duration(args.Duration * float64(time.Second)))
		w.Close()
		a.Quit()
	}()

	w.ShowAndRun()
}

type transparentFyneTheme struct{}

func (m transparentFyneTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m transparentFyneTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m transparentFyneTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m transparentFyneTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func drawRectangle(a *args) image.Image {
	dc := gg.NewContext(a.ScreenWidth, a.ScreenHeight)
	// Took me ages to guess at /2 here, but then it worked so :this_is_fine:...
	// If you're reading this, and you know why, please let me know as it's painful
	// to have fixed something but have no idea how you actually did it!!
	dc.DrawRectangle(a.X/2, a.Y/2, a.Width/2, a.Height/2)
	dc.SetColor(a.Colour)
	dc.SetLineWidth(a.LineWidth)
	dc.Stroke()

	return dc.Image()
}

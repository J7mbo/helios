package helios

import _ "embed"

type App struct {
	config *Config
	screen *Screen
}

//go:embed highlighter/highlighter.bin
var highlighterBinary []byte

func NewApplication(config *Config) *App {
	if config == nil {
		config, _ = NewConfig()
	}

	return &App{config: config, screen: NewScreen(NewHighlighter(highlighterBinary))}
}

func (a *App) GetScreen() *Screen {
	return a.screen
}

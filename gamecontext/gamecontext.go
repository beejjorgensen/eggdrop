package gamecontext

import (
	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/veandco/go-sdl2/sdl"
)

type gameContext struct {
	MainWindow  *sdl.Window
	MainSurface *sdl.Surface
	PixelFormat *sdl.PixelFormat

	GameManager *gamemanager.GameManager
}

// GContext holds the global game state
var GContext gameContext

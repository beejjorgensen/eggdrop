package gamecontext

import "github.com/veandco/go-sdl2/sdl"

type gameContext struct {
	MainWindow  *sdl.Window
	MainSurface *sdl.Surface
	PixelFormat *sdl.PixelFormat
}

// GContext holds the global game state
var GContext gameContext

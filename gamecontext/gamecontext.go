// Package gamecontext holds global, shared information about the system.
package gamecontext

import "github.com/veandco/go-sdl2/sdl"

type gameContext struct {
	MainWindow      *sdl.Window
	MainSurface     *sdl.Surface
	PixelFormat     *sdl.PixelFormat
	PixelFormatEnum uint32

	WindowWidth, WindowHeight int32
}

// GContext holds the global game state
var GContext = &gameContext{}

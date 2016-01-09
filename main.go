// Package main sets up the basic SDL window, makes the GameManager, and sets up
// the main event loop.
package main

import (
	"fmt"
	"runtime"

	"github.com/beejjorgensen/eggdrop/gamecontext"
	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/beejjorgensen/eggdrop/introstate"
	"github.com/beejjorgensen/eggdrop/playstate"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

func init() {
	// SDL wants to be on the main thread, OR ELSE!
	runtime.LockOSThread()
}

func sdlInit() {

	sdl.Init(sdl.INIT_EVERYTHING)

	// init the image subsystem
	imgFlags := img.INIT_PNG | img.INIT_JPG

	if imgFlags != img.Init(imgFlags) {
		panic(fmt.Sprintf("Error initializing img: %v\n", img.GetError()))
	}

	// init the TTF subsystem
	if err := ttf.Init(); err != nil {
		panic(fmt.Sprintf("Error initializing ttf: %v\n", err))
	}
}

func createMainWindow() {
	var err error
	var mainWindow *sdl.Window

	mainWindow, err = sdl.CreateWindow("Eggdrop!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	gamecontext.GContext.MainSurface, err = mainWindow.GetSurface()
	if err != nil {
		panic(err)
	}

	var winPixelFormat uint32

	winPixelFormat, err = mainWindow.GetPixelFormat()
	if err != nil {
		panic(err)
	}

	gamecontext.GContext.PixelFormat, err = sdl.AllocFormat(uint(winPixelFormat)) // TODO why the cast? Seems to work?
	if err != nil {
		panic(err)
	}

	gamecontext.GContext.MainWindow = mainWindow
}

func main() {
	sdlInit()

	gm := gamemanager.GGameManager

	createMainWindow()
	defer gamecontext.GContext.MainWindow.Destroy()

	intro := &introstate.IntroState{}
	play := &playstate.PlayState{}

	gm.RegisterMode(gamemanager.GameModeIntro, intro)
	gm.RegisterMode(gamemanager.GameModePlay, play)

	done := false

	// Keep handy for use in the loop
	mainWindow := gamecontext.GContext.MainWindow
	mainWindowSurface := gamecontext.GContext.MainSurface

	gm.SetMode(gamemanager.GameModeIntro)

	gm.SetEventMode(gamemanager.GameManagerEventDriven)

	for done == false {
		event := gm.GetNextEvent()
		for ; event != nil; event = sdl.PollEvent() {
			done = done || gm.HandleEvent(&event)

			switch event := event.(type) {
			case *sdl.WindowEvent:
				//fmt.Printf("Window: %#v\n", event)
				//fmt.Printf("Event: %t %t\n", event.Event, sdl.WINDOWEVENT_CLOSE)
				if event.Event == sdl.WINDOWEVENT_CLOSE {
					done = done || true
				}
			}
		}

		gm.Render(mainWindowSurface)
		mainWindow.UpdateSurface()

		gm.DelayToNextFrame()
	}

	sdl.Quit()
}

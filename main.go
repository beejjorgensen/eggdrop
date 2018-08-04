// Eggdrop is an SDL game written for the purposes of learning Go.
package main

import (
	"fmt"
	"runtime"

	"github.com/beejjorgensen/eggdrop/gamecontext"
	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/beejjorgensen/eggdrop/introstate"
	"github.com/beejjorgensen/eggdrop/playstate"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/ttf"
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

	gc := gamecontext.GContext
	gc.WindowWidth = 800
	gc.WindowHeight = 600

	mainWindow, err = sdl.CreateWindow("Eggdrop!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, gc.WindowWidth, gc.WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	gc.MainSurface, err = mainWindow.GetSurface()
	if err != nil {
		panic(err)
	}

	gc.PixelFormatEnum, err = mainWindow.GetPixelFormat()
	if err != nil {
		panic(err)
	}

	gc.PixelFormat, err = sdl.AllocFormat(uint(gc.PixelFormatEnum)) // TODO why the cast? Seems to work?
	if err != nil {
		panic(err)
	}

	gc.MainWindow = mainWindow
}

func main() {
	sdlInit()

	gm := gamemanager.GGameManager
	gc := gamecontext.GContext

	createMainWindow()
	defer gc.MainWindow.Destroy()

	intro := &introstate.IntroState{}
	play := &playstate.PlayState{}

	gm.RegisterMode(gamemanager.GameModeIntro, intro)
	gm.RegisterMode(gamemanager.GameModePlay, play)

	done := false

	// Keep handy for use in the loop
	mainWindow := gc.MainWindow
	mainWindowSurface := gc.MainSurface

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

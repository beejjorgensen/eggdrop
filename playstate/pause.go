package playstate

import (
	"fmt"

	"github.com/beejjorgensen/eggdrop/gamecontext"
	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/beejjorgensen/eggdrop/menu"
	"github.com/beejjorgensen/eggdrop/scenegraph"
	"github.com/beejjorgensen/eggdrop/util"
	"github.com/veandco/go-sdl2/sdl"
)

// Construct the pause menu
func (ps *PlayState) buildPauseMenu() {
	am := ps.assetManager // asset manager

	mainW := gamecontext.GContext.MainSurface.W
	mainH := gamecontext.GContext.MainSurface.H

	// Build pause menu shade background
	pauseBGSurface, err := util.MakeFillSurfaceAlpha(mainW, mainH, 0, 0, 0, 127)
	if err != nil {
		panic(fmt.Sprintf("Pause bgSurface: %v", err))
	}

	ps.pauseMenuEntity = scenegraph.NewEntity(pauseBGSurface)
	ps.pauseMenuEntity.Visible = false

	// Build pause menu
	mColor := ps.fontNormalColor
	mHiColor := ps.fontHighlightColor

	menuItems := []menu.Item{
		{AssetFontID: "menuFont", Text: "Return to Game", Color: mColor, HiColor: mHiColor},
		{AssetFontID: "menuFont", Text: "Main Menu", Color: mColor, HiColor: mHiColor},
	}

	ps.menu = menu.New(am, "playMenu", menuItems, 60, menu.MenuJustifyCenter)

	ps.menu.RootEntity.Y = 200

	ps.pauseMenuEntity.AddChild(ps.menu.RootEntity)
}

// pause pauses or unpauses the game
func (ps *PlayState) pause(paused bool) {
	gm := gamemanager.GGameManager

	if paused == ps.paused {
		return
	}

	if paused {
		// Show pause menu
		ps.menu.SetSelected(0)
		ps.pauseMenuEntity.Visible = true

		// Set to Event Driven
		gm.SetEventMode(gamemanager.GameManagerEventDriven)
	} else {
		// Hide pause menu
		ps.pauseMenuEntity.Visible = false

		// Set to Poll Driven
		gm.SetEventMode(gamemanager.GameManagerPollDriven)
	}
	ps.paused = paused
}

// handleEventPaused deals with paused events in the paused state
func (ps *PlayState) handleEventPaused(event *sdl.Event) bool {
	switch event := (*event).(type) {
	case *sdl.KeyboardEvent:
		//fmt.Printf("Key: %#v\n", event)
		switch event.Keysym.Sym {

		case sdl.K_ESCAPE:
			ps.pause(false)

		case sdl.K_DOWN:
			ps.menu.SelectNext()

		case sdl.K_UP:
			ps.menu.SelectPrev()

		case sdl.K_RETURN:
			if ps.handleMenuItem(ps.menu.GetSelected()) {
				return true // exit
			}
		}

	case *sdl.MouseMotionEvent:
		ps.menu.SelectByMouseY(event.Y)

	case *sdl.MouseButtonEvent:
		if event.Type == sdl.MOUSEBUTTONDOWN {
			ps.menu.SelectByMouseClickY(event.Y)

			clicked := ps.menu.GetClicked()
			if clicked >= 0 {
				if ps.handleMenuItem(clicked) {
					return true // exit
				}
			}
		}
	}

	return false
}

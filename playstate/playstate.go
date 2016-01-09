// Package playstate is responsible for assets and behavior of the actual game
// itself. This is where the action is. (Also handles the pause menu in the
// game.)
package playstate

import (
	"fmt"

	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/beejjorgensen/eggdrop/menu"
	"github.com/beejjorgensen/eggdrop/util"

	"github.com/beejjorgensen/eggdrop/assetmanager"
	"github.com/beejjorgensen/eggdrop/gamecontext"
	"github.com/beejjorgensen/eggdrop/scenegraph"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	assetMenuFontID = iota
	assetMenuBaseID = 0x10000
)

// PlayState holds information about the main game and pause menu
type PlayState struct {
	assetManager                *assetmanager.AssetManager
	rootEntity, pauseMenuEntity *scenegraph.Entity

	fontNormalColor, fontHighlightColor sdl.Color
	bgColor                             uint32

	paused bool

	menu *menu.Menu
}

// Init initializes this gamestate
func (ps *PlayState) Init() {
	// Create colors
	ps.bgColor = sdl.MapRGB(gamecontext.GContext.PixelFormat, 60, 60, 160)
	ps.fontNormalColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	ps.fontHighlightColor = sdl.Color{R: 255, G: 255, B: 0, A: 255}

	ps.assetManager = assetmanager.New()

	ps.loadAssets()
	ps.buildScene()
}

// buildScene constructs the necessary elements for the scene
func (ps *PlayState) buildScene() {
	mainW := gamecontext.GContext.MainSurface.W
	mainH := gamecontext.GContext.MainSurface.H

	rootEntity := scenegraph.NewEntity(nil)
	rootEntity.W = mainW
	rootEntity.H = mainH

	ps.buildPauseMenu()

	util.CenterEntityInParent(ps.menu.RootEntity, rootEntity)
	rootEntity.AddChild(ps.pauseMenuEntity)

	ps.rootEntity = rootEntity
}

// loadAssets loads this state's assets
func (ps *PlayState) loadAssets() {
	am := ps.assetManager // asset manager
	var err error

	if err = am.LoadFont(assetMenuFontID, "assets/Osborne1.ttf", 40); err != nil {
		panic(fmt.Sprintf("Playstate load font: %v", err))
	}
}

// handleMenuItem does the right thing with a selected menu item
func (ps *PlayState) handleMenuItem(i int) bool {
	switch i {
	case 0: // Continue
		ps.pause(false)
	case 1: // Quit
		// back to introstate
		gamemanager.GGameManager.SetMode(gamemanager.GameModeIntro)
	}

	return false
}

// handleEventPlaying deals with paused events in the play state
func (ps *PlayState) handleEventPlaying(event *sdl.Event) bool {
	switch event := (*event).(type) {
	case *sdl.KeyDownEvent:
		//fmt.Printf("Key: %#v\n", event)
		switch event.Keysym.Sym {

		case sdl.K_ESCAPE, sdl.K_p:
			ps.pause(true)
		}
	}

	return false
}

// HandleEvent handles SDL events for the intro state
func (ps *PlayState) HandleEvent(event *sdl.Event) bool {
	if ps.paused {
		return ps.handleEventPaused(event)
	}

	return ps.handleEventPlaying(event)
}

// Render renders the intro state
func (ps *PlayState) Render(mainWindowSurface *sdl.Surface) {
	rootEntity := ps.rootEntity

	mainWindowSurface.FillRect(nil, ps.bgColor)
	rootEntity.Render(mainWindowSurface)
}

// WillShow is called just before this state begins
func (ps *PlayState) WillShow() {
	// call this to move on to the next transition state
	gamemanager.GGameManager.WillShowComplete()
}

// WillHide is called just before this state ends
func (ps *PlayState) WillHide() {
}

// DidShow is called just after this statebegins
func (ps *PlayState) DidShow() {
}

// DidHide is called just after this state ends
func (ps *PlayState) DidHide() {
}

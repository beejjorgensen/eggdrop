// Package playstate is responsible for assets and behavior of the actual game
// itself. This is where the action is. (Also handles the pause menu in the
// game.)
package playstate

import (
	"fmt"

	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/beejjorgensen/eggdrop/menu"

	"github.com/beejjorgensen/eggdrop/assetmanager"
	"github.com/beejjorgensen/eggdrop/gamecontext"
	"github.com/beejjorgensen/eggdrop/scenegraph"
	"github.com/veandco/go-sdl2/sdl"
)

// PlayState holds information about the main game and pause menu
type PlayState struct {
	assetManager                *assetmanager.AssetManager
	rootEntity, pauseMenuEntity *scenegraph.Entity

	fontNormalColor, fontHighlightColor sdl.Color
	bgColor                             uint32

	paused bool

	menu *menu.Menu

	entityByID map[string]*scenegraph.Entity

	nestEntity *scenegraph.Entity
}

// Init initializes this gamestate
func (ps *PlayState) Init() {
	// Create colors
	ps.bgColor = sdl.MapRGB(gamecontext.GContext.PixelFormat, 133, 187, 234)
	ps.fontNormalColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	ps.fontHighlightColor = sdl.Color{R: 255, G: 255, B: 0, A: 255}

	ps.assetManager = assetmanager.New()
	ps.assetManager.SetOuterSurface(gamecontext.GContext.MainSurface)

	err := ps.assetManager.LoadJSON("playassets.json")
	if err != nil {
		panic(fmt.Sprintf("playassets.json: %v", err))
	}
	ps.buildScene()
}

// buildScene constructs the necessary elements for the scene
func (ps *PlayState) buildScene() {
	var err error

	am := ps.assetManager // asset manager

	ps.rootEntity, err = scenegraph.LoadJSON(am, "playgraph.json", nil)
	if err != nil {
		panic(fmt.Sprintf("playgraph.json: %v", err))
	}

	//ps.nestEntity = ps.entityByID["nest"]
	ps.nestEntity = ps.rootEntity.SearchByID("nest")

	// Pause menu stuff
	ps.buildPauseMenu()
	scenegraph.CenterEntityInParent(ps.menu.RootEntity, ps.rootEntity)
	ps.rootEntity.AddChild(ps.pauseMenuEntity)
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

// positionNest positions and clamps the nest
func (ps *PlayState) positionNest(x int32) {
	w := ps.nestEntity.W
	x -= w / 2 // center

	if x < 0 {
		x = 0
	}

	maxX := gamecontext.GContext.WindowWidth - w
	if x > maxX {
		x = maxX
	}

	ps.nestEntity.X = x
}

// handleEventPlaying deals with events in the play state
func (ps *PlayState) handleEventPlaying(event *sdl.Event) bool {
	switch event := (*event).(type) {
	case *sdl.KeyDownEvent:
		//fmt.Printf("Key: %#v\n", event)
		switch event.Keysym.Sym {

		case sdl.K_ESCAPE, sdl.K_p:
			ps.pause(true)
		}
	case *sdl.MouseMotionEvent:
		ps.positionNest(event.X)
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
	ps.pause(false)

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

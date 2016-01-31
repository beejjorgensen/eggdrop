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
	stateInterludeDuration       = 250  // ms
	stateInterludeDurationLevel1 = 1000 // ms
)

const (
	stateInterlude = iota
	stateAction
)

type stateInfo struct {
	state     int
	startTime uint32
}

// PlayState holds information about the main game and pause menu
type PlayState struct {
	assetManager                *assetmanager.AssetManager
	rootEntity, pauseMenuEntity *scenegraph.Entity

	fontNormalColor, fontHighlightColor sdl.Color
	bgColor                             uint32

	paused bool

	menu *menu.Menu

	nestEntity          *scenegraph.Entity
	chixEntity          *scenegraph.Entity
	chixLeftEntity      *scenegraph.Entity
	chixRightEntity     *scenegraph.Entity
	interludeTextEntity *scenegraph.Entity
	eggContainer        *scenegraph.Entity
	chixLegEntity       []*scenegraph.Entity

	eggTimeSinceLaunch uint32
	eggLaunchDelay     uint32

	state stateInfo
	level int

	chix chixInfo
}

// Init initializes this gamestate
func (ps *PlayState) Init() {
	ps.initChix()
	ps.initEggs()

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

	ps.nestEntity = ps.rootEntity.SearchByID("nest")
	ps.chixEntity = ps.rootEntity.SearchByID("chicken")
	ps.chixLeftEntity = ps.rootEntity.SearchByID("chickenLeftContainer")
	ps.chixRightEntity = ps.rootEntity.SearchByID("chickenRightContainer")
	ps.chixLegEntity = []*scenegraph.Entity{
		ps.rootEntity.SearchByID("chickenLeftLegs0"),
		ps.rootEntity.SearchByID("chickenLeftLegs1"),
		ps.rootEntity.SearchByID("chickenRightLegs0"),
		ps.rootEntity.SearchByID("chickenRightLegs1"),
	}
	ps.interludeTextEntity = ps.rootEntity.SearchByID("interludeText")
	ps.eggContainer = ps.rootEntity.SearchByID("eggContainer")

	// This is hackish, but we need to know the width of the chicken, and
	// the chicken parent node is sizeless. So we copy the size from one of
	// the children. This should probably be an option in the JSON reader.
	ps.chixEntity.W = ps.rootEntity.SearchByID("chickenLeftImage").W

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
		if ps.state.state == stateAction {
			ps.positionNest(event.X)
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

// updateState watches for internal state changes
func (ps *PlayState) updateState() {
	curTime := sdl.GetTicks()
	diff := curTime - ps.state.startTime

	switch ps.state.state {
	case stateInterlude:
		var duration uint32
		switch ps.level {
		case 1:
			duration = stateInterludeDurationLevel1
		default:
			duration = stateInterludeDuration
		}
		if diff >= duration {
			ps.setState(stateAction)
		}
	case stateAction:
	}

}

// update handles the updating of entities, timers, time state changes, etc. per
// frame
func (ps *PlayState) update() {
	ps.updateState()

	switch ps.state.state {
	case stateAction:
		ps.updateChix()
		ps.updateEggs()
	}
}

// Render renders the intro state
func (ps *PlayState) Render(mainWindowSurface *sdl.Surface) {
	ps.update()

	mainWindowSurface.FillRect(nil, ps.bgColor)

	ps.interludeTextEntity.Visible = ps.state.state == stateInterlude

	ps.rootEntity.Render(mainWindowSurface)
}

// constructInterludeImage builds the "LEVEL X" image
func (ps *PlayState) constructInterludeImage() {
	surface, err := util.RenderText(ps.assetManager.Fonts["interludeFont"], fmt.Sprintf("LEVEL %d", ps.level), sdl.Color{R: 255, G: 255, B: 0, A: 255})

	if err != nil {
		panic(fmt.Sprintf("Error constructing interlude text: %v", err))
	}

	ps.interludeTextEntity.Surface = surface
	ps.interludeTextEntity.W = surface.W // hackishly make centering work on the next line
	scenegraph.CenterEntityInParent(ps.interludeTextEntity, ps.rootEntity)
}

// setState sets the current states and does timer management
func (ps *PlayState) setState(state int) {
	ps.state.state = state
	ps.state.startTime = sdl.GetTicks()
}

// WillShow is called just before this state begins
func (ps *PlayState) WillShow() {
	ps.resetChix()
	ps.pause(false)
	ps.level = 1
	ps.setState(stateInterlude)
	ps.constructInterludeImage()

	// call this to move on to the next transition state
	gamemanager.GGameManager.WillShowComplete()
}

// WillHide is called just before this state ends
func (ps *PlayState) WillHide() {
}

// DidShow is called just after this statebegins
func (ps *PlayState) DidShow() {
	gamemanager.GGameManager.SetEventMode(gamemanager.GameManagerPollDriven)
}

// DidHide is called just after this state ends
func (ps *PlayState) DidHide() {
}

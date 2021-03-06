// Package introstate controls assets and behavior for the initial intro state
// of the game.
package introstate

import (
	"fmt"

	"github.com/beejjorgensen/eggdrop/assetmanager"
	"github.com/beejjorgensen/eggdrop/gamecontext"
	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/beejjorgensen/eggdrop/menu"
	"github.com/beejjorgensen/eggdrop/scenegraph"
	"github.com/veandco/go-sdl2/sdl"
)

// IntroState holds all information about the intro state
type IntroState struct {
	assetManager                        *assetmanager.AssetManager
	rootEntity                          *scenegraph.Entity
	bgColor                             uint32
	fontNormalColor, fontHighlightColor sdl.Color
	menu                                *menu.Menu
}

// Init initializes this gamestate
func (is *IntroState) Init() {
	// Create colors
	is.bgColor = sdl.MapRGB(gamecontext.GContext.PixelFormat, 60, 160, 60)
	is.fontNormalColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	is.fontHighlightColor = sdl.Color{R: 255, G: 255, B: 0, A: 255}

	is.assetManager = assetmanager.New()

	err := is.assetManager.LoadJSON("introassets.json")
	if err != nil {
		panic(fmt.Sprintf("introassets.json: %v", err))
	}

	is.buildScene()
}

func (is *IntroState) buildScene() {
	am := is.assetManager // asset manager

	rootEntity := scenegraph.NewEntity(nil)
	rootEntity.W = gamecontext.GContext.MainSurface.W
	rootEntity.H = gamecontext.GContext.MainSurface.H

	titleEntity := scenegraph.NewEntity(am.Surfaces["titleText"])

	mColor := is.fontNormalColor
	mHiColor := is.fontHighlightColor

	menuItems := []menu.Item{
		{AssetFontID: "menuFont", Text: "Play!", Color: mColor, HiColor: mHiColor},
		{AssetFontID: "menuFont", Text: "Quit", Color: mColor, HiColor: mHiColor},
	}

	is.menu = menu.New(am, "introMenu", menuItems, 60, menu.MenuJustifyCenter)

	scenegraph.CenterEntityInParent(is.menu.RootEntity, rootEntity)
	is.menu.RootEntity.Y = 200

	rootEntity.AddChild(titleEntity, is.menu.RootEntity)

	is.rootEntity = rootEntity

	// position title
	scenegraph.CenterEntityInSurface(titleEntity, gamecontext.GContext.MainSurface)
	titleEntity.Y = 40

}

// handleMenuItem does the right thing with a selected menu item
func (is *IntroState) handleMenuItem(i int) bool {
	switch i {
	case 0: // Play!
		gamemanager.GGameManager.SetMode(gamemanager.GameModePlay)
	case 1: // Quit
		return true // exit game
	}

	return false
}

// HandleEvent handles SDL events for the intro state
func (is *IntroState) HandleEvent(event *sdl.Event) bool {

	switch event := (*event).(type) {
	case *sdl.KeyboardEvent:
		//fmt.Printf("Key: %#v\n", event)
		switch event.Keysym.Sym {

		case sdl.K_ESCAPE:
			return true // exit game

		case sdl.K_DOWN:
			is.menu.SelectNext()

		case sdl.K_UP:
			is.menu.SelectPrev()

		case sdl.K_RETURN:
			if is.handleMenuItem(is.menu.GetSelected()) {
				return true // exit
			}
		}

	case *sdl.MouseMotionEvent:
		is.menu.SelectByMouseY(event.Y)

	case *sdl.MouseButtonEvent:
		if event.Type == sdl.MOUSEBUTTONDOWN {
			is.menu.SelectByMouseClickY(event.Y)

			clicked := is.menu.GetClicked()
			if clicked >= 0 {
				if is.handleMenuItem(clicked) {
					return true // exit
				}
			}
		}
	}

	return false
}

// Render renders the intro state
func (is *IntroState) Render(mainWindowSurface *sdl.Surface) {
	rootEntity := is.rootEntity

	mainWindowSurface.FillRect(nil, is.bgColor)
	rootEntity.Render(mainWindowSurface)
}

// WillShow is called just before this state begins
func (is *IntroState) WillShow() {
	// call this to move on to the next transition state
	gamemanager.GGameManager.WillShowComplete()
}

// WillHide is called just before this state ends
func (is *IntroState) WillHide() {
}

// DidShow is called just after this statebegins
func (is *IntroState) DidShow() {
	gamemanager.GGameManager.SetEventMode(gamemanager.GameManagerEventDriven)
}

// DidHide is called just after this state ends
func (is *IntroState) DidHide() {
}

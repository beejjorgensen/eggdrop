// Package gamemanager tracks the current state of the game. States are
// registered in advance, and satisfy the GameMode interface. When states are
// changed, the old and new states are called with appropriate notifcations so
// they can prepare for the change.
//
// Also can execute the frame delay, and poll for SDL events in a number of ways
// (Event, EventWithTimeout, Poll).
package gamemanager

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// GameMode is methods for handling game events and state changes
type GameMode interface {
	Init()
	Render(*sdl.Surface)
	HandleEvent(*sdl.Event) bool
	WillShow()
	DidShow()
	WillHide()
	DidHide()
}

// Event modes for eventMode
const (
	GameManagerEventDriven = iota
	GameManagerEventTimeoutDriven
	GameManagerPollDriven
)

// GameManager manages the main game states
type GameManager struct {
	currentModeID int
	nextModeID    int
	modeMap       map[int]GameMode

	FrameDelay   uint32 // ms
	EventTimeout int
	eventMode    int

	prevFrameTime uint32
}

// GGameManager is the global game manager
var GGameManager = New()

// New creates a new initialized GameManager
func New() *GameManager {
	return &GameManager{
		currentModeID: -1,
		modeMap:       make(map[int]GameMode),
		FrameDelay:    1000 / 60,
		EventTimeout:  1000 / 60,
	}
}

// RegisterMode registers a new main game mode
func (g *GameManager) RegisterMode(id int, gm GameMode) {
	g.modeMap[id] = gm
	gm.Init()
}

// SetMode sets the current game mode to the specified value
func (g *GameManager) SetMode(id int) {
	g.nextModeID = id

	if g.currentModeID >= 0 {
		g.modeMap[g.currentModeID].WillHide()
	}
	g.modeMap[id].WillShow()
}

// WillShowComplete tells the GameManager it's time to go to the next stage
// of the SetMode sequence
func (g *GameManager) WillShowComplete() {
	if g.currentModeID >= 0 {
		g.modeMap[g.currentModeID].DidHide()
	}
	g.modeMap[g.nextModeID].DidShow()

	g.currentModeID = g.nextModeID
}

// HandleEvent forwards to the event handler for the current GameMode
func (g *GameManager) HandleEvent(event *sdl.Event) bool {
	return g.modeMap[g.currentModeID].HandleEvent(event)
}

// Render forwards to the renderer for the current GameMode
func (g *GameManager) Render(surface *sdl.Surface) {
	g.modeMap[g.currentModeID].Render(surface)
}

// DelayToNextFrame waits until it's time to do the next event/render loop
func (g *GameManager) DelayToNextFrame() {
	curTime := sdl.GetTicks()

	if g.prevFrameTime == 0 {
		if curTime >= g.FrameDelay {
			g.prevFrameTime = curTime - g.FrameDelay
		}
	}

	diff := curTime - g.prevFrameTime

	if g.FrameDelay > diff {
		frameDelayUnder := g.FrameDelay - diff
		// we have not yet exceeded one frame, so we need to sleep
		//fmt.Printf("Under: %d %d %d %d\n", curTime, g.prevFrameTime, diff, frameDelayUnder)
		sdl.Delay(frameDelayUnder)
	} else {
		//frameDelayOver := diff - g.FrameDelay
		//fmt.Printf("Over: %d %d %d %d\n", curTime, g.prevFrameTime, diff, frameDelayOver)
		// we have exceeded one frame, so no sleep
		// TODO sleep less in the future to make up for it?
	}

	g.prevFrameTime = curTime
}

// SetEventMode sets the event handler mode to polling or event-based
func (g *GameManager) SetEventMode(eventMode int) {
	g.eventMode = eventMode

	// Push an empty event to kick out of WaitEvent() or WaitEventTimeout()
	sdl.PushEvent(&sdl.UserEvent{Type: 0})
}

// GetNextEvent returns the next sdl.Event depending on the eventMode
func (g *GameManager) GetNextEvent() sdl.Event {
	switch g.eventMode {
	case GameManagerEventDriven:
		return sdl.WaitEvent()
	case GameManagerEventTimeoutDriven:
		return sdl.WaitEventTimeout(g.EventTimeout)
	case GameManagerPollDriven:
		return sdl.PollEvent()
	}

	panic(fmt.Sprintf("GetNextEvent: unknown event mode: %d", g.eventMode))
}

package playstate

import (
	"math"
	"math/rand"

	"github.com/beejjorgensen/eggdrop/gamemanager"
)

const (
	chixInitAngleSpeed     = 0.001 // Higher is faster motion
	chixInitAngleSpeedMult = 1.2   // How much higher to go each level, multiplier
)

// chixInfo holds the chicken state
type chixInfo struct {
	pos, prevPos int32
	angleSpeed   float64
	angle        float64
}

// initChix creates an empty chicken state
func (ps *PlayState) initChix() {
	ps.chix = chixInfo{}
}

// resetChix restores the chicken to the start state
func (ps *PlayState) resetChix() {
	ps.chix.pos = 0
	ps.chix.prevPos = 0
	ps.chix.angleSpeed = chixInitAngleSpeed

	// The chicken will be centered at multiples of 2π, so we choose one of 100000
	// of those arbitrarily
	ps.chix.angle = float64(rand.Int31n(100000)*2) * math.Pi
}

// updateChix updates and positions the chicken
func (ps *PlayState) updateChix() {
	// Use the optimal frame delay to update positions
	frameDelay := gamemanager.GGameManager.FrameDelay

	β := ps.chix.angle + float64(frameDelay)*ps.chix.angleSpeed

	// Compute angle [-1..1]
	// I just made up these numbers. They probably don't interfere nearly as much
	// as I hope they do
	pos := (math.Sin(β) + math.Sin(β*2.5) + math.Sin(β*4.7)) / 3.0

	// Remap to screen
	halfScreen := ps.rootEntity.W / 2
	posInt := int32(pos*float64(halfScreen)*0.85) + halfScreen // range over 85% of screen width

	ps.chix.pos = posInt - ps.chixEntity.W/2

	movingRight := ps.chix.pos > ps.chix.prevPos
	ps.chixLeftEntity.Visible = !movingRight
	ps.chixRightEntity.Visible = movingRight

	ps.chixEntity.X = ps.chix.pos
	ps.chix.prevPos = ps.chix.pos

	ps.chix.angle = β
}

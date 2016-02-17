package playstate

import (
	"github.com/beejjorgensen/eggdrop/gamemanager"
	"github.com/beejjorgensen/eggdrop/scenegraph"
)

const (
	eggStartingY = 50  // pixels
	eggSplatY    = 570 // pixels

	eggSpeed0        = 250 // pixels per second
	eggSpeedPerLevel = 100 // pixels per second, how much to increase per level

	eggDelay0        = 500  // ms
	eggDelayPerLevel = 0.85 // multiple

	eggLaunchXOffset = 80 // px
)

// initEggs initializes the egg subsystem. It sounds so important when I call it
// that.
func (ps *PlayState) initEggs() {
	ps.eggLaunchDelay = eggDelay0
}

// resetEggs hides all the eggs
func (ps *PlayState) resetEggs() {
	for _, egg := range ps.eggContainer.Children {
		egg.Visible = false
	}
}

// newEgg creates a new egg and adds it to the egg container
func (ps *PlayState) newEgg() *scenegraph.Entity {
	egg := scenegraph.NewEntity(ps.assetManager.Surfaces["eggImage"])
	egg.Visible = false
	ps.eggContainer.AddChild(egg)

	return egg
}

// getEgg returns a ready-to-use egg, creating it if necessary. This is an O(N)
// process, sorry.
func (ps *PlayState) getEgg() *scenegraph.Entity {
	var readyEgg *scenegraph.Entity

	for _, egg := range ps.eggContainer.Children {
		if !egg.Visible {
			readyEgg = egg
			break
		}
	}

	// If we didn't find one, make a new one
	if readyEgg == nil {
		readyEgg = ps.newEgg()
	}

	return readyEgg
}

// launchEgg brings a new egg into existence
func (ps *PlayState) launchEgg() {
	egg := ps.getEgg()

	egg.Y = eggStartingY

	var offset int32

	if ps.chix.Direction > 0 { // right
		offset = eggLaunchXOffset
	} else { // left
		offset = ps.chixEntity.W - eggLaunchXOffset
	}

	egg.X = ps.chixEntity.X + offset

	egg.Visible = true
}

// updateEggs animates eggs to their new position
func (ps *PlayState) updateEggs() {
	frameDelay := gamemanager.GGameManager.FrameDelay

	// Drop new eggs
	ps.eggTimeSinceLaunch += frameDelay
	if ps.eggTimeSinceLaunch > ps.eggLaunchDelay {
		ps.launchEgg()
		ps.eggTimeSinceLaunch = 0
	}

	// Animate eggs

	speed := eggSpeed0 + ps.level*eggSpeedPerLevel // px per second
	dY := int32(speed * int(frameDelay) / 1000)

	// Right now it just iterates through all eggs and animates the visible ones.
	// This could be improved for efficiency.
	for _, egg := range ps.eggContainer.Children {
		if egg.Visible {
			egg.MoveTo(egg.X, egg.Y+dY)

			if egg.Y > eggSplatY {
				egg.Visible = false
			}
		}
	}
}

// testEggCollision looks for collisions with the eggs and nest
func (ps *PlayState) testEggCollision() {
	for _, egg := range ps.eggContainer.Children {
		if egg.Visible {
			if egg.MoveAABB.TestCollision(&ps.nestEntity.MoveAABB) {
				egg.Visible = false
				//fmt.Printf("Hit!\n%#v\n%#v\n", egg.MoveAABB, ps.nestEntity.MoveAABB)
			}
		}
	}

}

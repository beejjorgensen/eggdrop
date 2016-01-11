package playstate

import (
	"fmt"

	"github.com/beejjorgensen/eggdrop/util"
)

func (ps *PlayState) loadChicken() {
	am := ps.assetManager // asset manager

	chixLeftSurface, err := am.LoadSurface(assetChickenLeftID, "assets/chicken.png")
	if err != nil {
		panic(fmt.Sprintf("chicken.png: %v", err))
	}

	chixRightSurface, err := util.SurfaceFlipH(chixLeftSurface)
	if err != nil {
		panic(fmt.Sprintf("chixRightSurface: %v", err))
	}
	am.AddSurface(assetChickenRightID, chixRightSurface)
}

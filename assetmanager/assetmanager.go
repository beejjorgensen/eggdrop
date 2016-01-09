// Package assetmanager loads assets (images and fonts) so they can be referred
// to later by ID.
package assetmanager

import (
	"github.com/beejjorgensen/eggdrop/util"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

// AssetManager holds surfaces and other asset information
type AssetManager struct {
	Surfaces map[int]*sdl.Surface
	Fonts    map[int]*ttf.Font
}

// New creates and initializes a new AssetManager
func New() *AssetManager {
	return &AssetManager{
		Surfaces: make(map[int]*sdl.Surface),
		Fonts:    make(map[int]*ttf.Font),
	}
}

// LoadSurface loads and tracks a new image
func (am *AssetManager) LoadSurface(key int, fileName string) (surface *sdl.Surface, err error) {
	surface, err = img.Load(fileName)
	am.Surfaces[key] = surface

	return
}

// AddSurface tracks an existing surface
func (am *AssetManager) AddSurface(key int, surface *sdl.Surface) {
	am.Surfaces[key] = surface
}

// LoadFont loads and tracks a Font
func (am *AssetManager) LoadFont(key int, fileName string, size int) (err error) {
	am.Fonts[key], err = ttf.OpenFont(fileName, size)

	return
}

// RenderText is a helper function to generate and track a surface with some text on it
func (am *AssetManager) RenderText(surfaceKey, fontKey int, text string, color sdl.Color) (*sdl.Surface, error) {
	surface, err := util.RenderText(am.Fonts[fontKey], text, color)

	if err == nil {
		am.AddSurface(surfaceKey, surface)
	}

	return surface, err
}

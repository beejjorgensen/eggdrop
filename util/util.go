// Package util contains misc utility functions.
package util

import (
	"reflect"
	"unsafe"

	"github.com/beejjorgensen/eggdrop/scenegraph"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

// CenterEntityInParent centers an entity within its parent
func CenterEntityInParent(entity, parent *scenegraph.Entity) {
	entity.X = (parent.W - entity.W) / 2
}

// CenterEntityInSurface centers an entity within a surface
func CenterEntityInSurface(entity *scenegraph.Entity, surface *sdl.Surface) {
	entity.X = (surface.W - entity.W) / 2
}

// CenterEntity centers an entity within a given width
func CenterEntity(entity *scenegraph.Entity, width int32) {
	entity.X = (width - entity.W) / 2
}

// RightJustifyEntityInParent right justifies an entity within its parent
func RightJustifyEntityInParent(entity, parent *scenegraph.Entity) {
	entity.X = parent.W - entity.W
}

// RightJustifyEntityInSurface right justifies an entity within a surface
func RightJustifyEntityInSurface(entity *scenegraph.Entity, surface *sdl.Surface) {
	entity.X = surface.W - entity.W
}

// RightJustifyEntity right justifies an entity within a given width
func RightJustifyEntity(entity *scenegraph.Entity, width int32) {
	entity.X = width - entity.W
}

// RenderText is a helper function to generate a surface with some text on it
func RenderText(font *ttf.Font, text string, color sdl.Color) (*sdl.Surface, error) {
	surface, err := font.RenderUTF8_Solid(text, color)

	return surface, err
}

// MakeFillSurfaceAlpha makes a new rectangular surface and fills with with RGBA
func MakeFillSurfaceAlpha(width, height int32, red, green, blue, alpha uint8) (*sdl.Surface, error) {
	pf, err := sdl.AllocFormat(sdl.PIXELFORMAT_RGBA8888)
	if err != nil {
		return nil, err
	}
	color := sdl.MapRGBA(pf, red, green, blue, alpha)
	surface, err := sdl.CreateRGBSurface(0, width, height, 32, pf.Rmask, pf.Gmask, pf.Bmask, pf.Amask)
	if err != nil {
		return nil, err
	}
	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: width, H: height}, color)

	return surface, nil
}

// MakeFillSurfaceConvertFormat makes a new rectangular surface and fills with with RGBA
func MakeFillSurfaceConvertFormat(width, height int32, red, green, blue, alpha uint8, pf uint32) (*sdl.Surface, error) {
	surface, err := MakeFillSurfaceAlpha(width, height, red, green, blue, alpha)
	if err != nil {
		return nil, err
	}

	return surface.ConvertFormat(pf, 0)
}

// MakeEmptySurfaceFrom makes a new empty surface the same size and format of the given surface.
// WARNING: does not handle palettized images!
func MakeEmptySurfaceFrom(src *sdl.Surface) (*sdl.Surface, error) {
	pf := src.Format

	return sdl.CreateRGBSurface(0, src.W, src.H, int32(pf.BitsPerPixel), pf.Rmask, pf.Gmask, pf.Bmask, pf.Amask)
}

// SurfaceFlipH produces a new surface with a horizontally-flipped version of the source surface.
func SurfaceFlipH(src *sdl.Surface) (*sdl.Surface, error) {
	pf := src.Format

	pixels := src.Pixels()
	pixelsDest := make([]byte, len(pixels))
	bytesPP := int32(pf.BytesPerPixel)
	bitsPP := int32(pf.BitsPerPixel)

	srcWPx := src.W * bytesPP

	for y := src.H - 1; y >= 0; y-- {
		rowStart := y * srcWPx
		for x := src.W - 1; x >= 0; x-- {
			si := rowStart + x*bytesPP
			di := rowStart + (src.W-x-1)*bytesPP

			copy(pixelsDest[di:], pixels[si:si+bytesPP])
		}
	}

	// We get the slice as a SliceHeader so we extract the .Data from it to
	// pass to CreateRGBSurfaceFrom():
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&pixelsDest))

	return sdl.CreateRGBSurfaceFrom(unsafe.Pointer(sh.Data), int(src.W), int(src.H), int(bitsPP), int(src.Pitch), pf.Rmask, pf.Gmask, pf.Bmask, pf.Amask)
}

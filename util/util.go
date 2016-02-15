// Package util contains misc utility functions.
package util

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

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

// surfaceManipulate is an internal function that creates a new surface, gets
// metadata, and then calls a passed-in function to actually do something to it,
// e.g. horizontally flip all the pixels.
//
// Note: this only supports manipulations that leave the Surface the same
// dimensions. It could potentially be generalized to handle manipulations that
// leave the Surface with the same number of pixels.
func surfaceManipulate(src *sdl.Surface, f func(src *sdl.Surface, srcWBytes, bytesPerPixel int32, srcPx, destPx []byte)) (*sdl.Surface, error) {
	pf := src.Format

	pixels := src.Pixels()
	bytesPP := int32(pf.BytesPerPixel)
	bitsPP := int32(pf.BitsPerPixel)

	newSurface, err := sdl.CreateRGBSurface(0, src.W, src.H, bitsPP, pf.Rmask, pf.Gmask, pf.Bmask, pf.Amask)

	if err != nil {
		return nil, err
	}

	pixelsDest := newSurface.Pixels()

	srcWBytes := int32(src.W * bytesPP)

	// Perform the manipulation
	f(src, srcWBytes, bytesPP, pixels, pixelsDest)

	return newSurface, nil
}

// SurfaceFlipH produces a new surface with a horizontally-flipped version of
// the source surface.
func SurfaceFlipH(src *sdl.Surface) (*sdl.Surface, error) {
	return surfaceManipulate(src, func(src *sdl.Surface, srcWBytes, bytesPerPixel int32, srcPx, destPx []byte) {
		for y := src.H - 1; y >= 0; y-- {
			rowStart := y * srcWBytes
			for x := src.W - 1; x >= 0; x-- {
				si := rowStart + x*bytesPerPixel
				di := rowStart + (src.W-x-1)*bytesPerPixel

				copy(destPx[di:], srcPx[si:si+bytesPerPixel])
			}
		}
	})
}

// SurfaceFlipV produces a new surface with a vertically-flipped version of the
// source surface.
func SurfaceFlipV(src *sdl.Surface) (*sdl.Surface, error) {
	return surfaceManipulate(src, func(src *sdl.Surface, srcWBytes, bytesPerPixel int32, srcPx, destPx []byte) {
		for y := src.H - 1; y >= 0; y-- {
			rowStart := y * srcWBytes
			dRowStart := (src.H - y - 1) * srcWBytes
			for x := src.W - 1; x >= 0; x-- {
				si := rowStart + x*bytesPerPixel
				di := dRowStart + x*bytesPerPixel

				copy(destPx[di:], srcPx[si:si+bytesPerPixel])
			}
		}
	})
}

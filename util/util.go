// Package util contains misc utility functions.
package util

import (
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

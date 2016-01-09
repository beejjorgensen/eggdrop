// Package scenegraph is a simple hierarchical scene graph, and declares an
// Entity type that associates an entity with a surface.
//
// The EntityTransform type is currently cheesy, but could eventually be
// replaced by something more appropriate (e.g. a matrix).
package scenegraph

import (
	"github.com/veandco/go-sdl2/sdl"
)

const initialChildrenCap = 5

// EntityTransform holds Entity transform information.
// Maybe someday this will grow up to be some kind of matrix.
type EntityTransform struct {
	X, Y, W, H int32
}

// Entity is a graphics entity with associated surface
type Entity struct {
	EntityTransform
	EntityToWorld, WorldToEntity EntityTransform
	Surface                      *sdl.Surface
	children                     []*Entity
	Visible                      bool
}

// NewEntity creates a new Entity for a given surface (or nil)
func NewEntity(surface *sdl.Surface) *Entity {
	e := &Entity{
		Surface:  surface,
		children: make([]*Entity, 0, initialChildrenCap),
		Visible:  true,
	}

	if surface != nil {
		e.W = surface.W
		e.H = surface.H
	}

	return e
}

// AddChild adds a child node to this entity
func (e *Entity) AddChild(childs ...*Entity) { // Tribute to Childs from The Thing
	for _, c := range childs {
		// TODO: Growth code could be optimized out of this loop
		currentCap := cap(e.children)
		if len(e.children) == currentCap {
			newChildren := make([]*Entity, currentCap*2)
			copy(newChildren, e.children)
			e.children = append(newChildren, c)
		} else {
			e.children = append(e.children, c)
		}
	}
}

// GetChild returns a specific child by index
func (e *Entity) GetChild(index int) *Entity {
	return e.children[index]
}

// Internal render call
func (e *Entity) renderRecursive(dest *sdl.Surface, t EntityTransform) {
	// If invisible, stop processing this subtree
	if !e.Visible {
		return
	}

	// save these for anyone who might want them later
	e.EntityToWorld.X = e.X + t.X
	e.EntityToWorld.Y = e.Y + t.Y
	e.WorldToEntity.X = -e.EntityToWorld.X // invert
	e.WorldToEntity.Y = -e.EntityToWorld.Y

	rect := sdl.Rect{X: e.X + t.X, Y: e.Y + t.Y}
	e.Surface.Blit(nil, dest, &rect)

	t.X += e.X
	t.Y += e.Y

	for _, c := range e.children {
		c.renderRecursive(dest, t)
	}
}

// Render renders a hierarchy to the given surface
func (e *Entity) Render(dest *sdl.Surface) {

	t := EntityTransform{0, 0, dest.W, dest.H}

	e.renderRecursive(dest, t)
}

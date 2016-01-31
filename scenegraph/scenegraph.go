// Package scenegraph is a simple hierarchical scene graph, and declares an
// Entity type that associates an entity with a surface.
//
// The EntityTransform type is currently cheesy, but could eventually be
// replaced by something more appropriate (e.g. a matrix).
package scenegraph

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/beejjorgensen/eggdrop/assetmanager"
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
	Children                     []*Entity
	Visible                      bool
	ID                           string
}

// NewEntity creates a new Entity for a given surface (or nil)
func NewEntity(surface *sdl.Surface) *Entity {
	e := &Entity{
		Surface:  surface,
		Children: make([]*Entity, 0, initialChildrenCap),
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
		currentCap := cap(e.Children)
		if len(e.Children) == currentCap {
			newChildren := make([]*Entity, currentCap, currentCap*2)
			copy(newChildren, e.Children)
			e.Children = append(newChildren, c)
		} else {
			e.Children = append(e.Children, c)
		}
	}
}

// GetChild returns a specific child by index
func (e *Entity) GetChild(index int) *Entity {
	return e.Children[index]
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

	for _, c := range e.Children {
		c.renderRecursive(dest, t)
	}
}

// Render renders a hierarchy to the given surface
func (e *Entity) Render(dest *sdl.Surface) {

	t := EntityTransform{0, 0, dest.W, dest.H}

	e.renderRecursive(dest, t)
}

// parseJSONPos
// loadJSONRecursive runs down the hierarchy in the JSON file
func loadJSONRecursive(am *assetmanager.AssetManager, node map[string]interface{}, filename string, entityByID map[string]*Entity) *Entity {
	var id string
	var idOK, ok bool

	var childEntity *Entity

	entity := NewEntity(nil)

	// Do this first to help with error messaging
	id, idOK = node["Id"].(string)
	if idOK {
		entity.ID = id
		if entityByID != nil {
			entityByID[id] = entity
		}
	} else {
		id = "[No ID]"
	}

	// Do these first since X and Y might depend on them:
	_, ok = node["W"].(string)
	if ok {
		entity.W = am.ParseDimension(filename, id, "W", node)
	}
	_, ok = node["H"].(string)
	if ok {
		entity.H = am.ParseDimension(filename, id, "H", node)
	}
	assetName, ok := node["Asset"].(string)
	if ok {
		entity.Surface = am.Surfaces[assetName]
		entity.W = entity.Surface.W
		entity.H = entity.Surface.H
	}

	// Now do the rest of the properties
	for k, v := range node {
		switch k {
		case "Id", "W", "H", "Asset":
			// do nothing; handled above
		case "X":
			entity.X = am.ParsePosition(filename, id, entity.W, "X", node)
		case "Y":
			entity.Y = am.ParsePosition(filename, id, entity.H, "Y", node)
		case "Visible":
			entity.Visible = v.(bool)
		case "Children":
			for _, child := range v.([]interface{}) {
				childEntity = loadJSONRecursive(am, child.(map[string]interface{}), filename, entityByID)
				entity.AddChild(childEntity)
			}
		default:
			panic(fmt.Sprintf("scenegraph.LoadJSON: unrecognized key: %s", k))
		}
	}

	return entity
}

// LoadJSON loads the scene from a JSON file
func LoadJSON(am *assetmanager.AssetManager, jsonFile string, entityByID map[string]*Entity) (*Entity, error) {
	jsonStr, err := ioutil.ReadFile(assetmanager.AssetPath(jsonFile))

	var jsonRoot map[string]interface{} // outer JSON is an object

	err = json.Unmarshal(jsonStr, &jsonRoot)
	if err != nil {
		return nil, err
	}

	sceneRoot, ok := jsonRoot["Scenegraph"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Can't find Scenegraph elemnt")
	}

	root := loadJSONRecursive(am, sceneRoot, jsonFile, entityByID)

	return root, nil
}

// SearchByID returns the entity with the matching ID, or nil. Slow. Pass
// entityByID map into LoadJSON for quicker results.
func (e *Entity) SearchByID(id string) *Entity {
	if e.ID == id {
		return e
	}

	for _, child := range e.Children {
		found := child.SearchByID(id)
		if found != nil {
			return found
		}
	}

	return nil
}

// CenterEntityInParent centers an entity within its parent
func CenterEntityInParent(entity, parent *Entity) {
	entity.X = (parent.W - entity.W) / 2
}

// CenterEntityInSurface centers an entity within a surface
func CenterEntityInSurface(entity *Entity, surface *sdl.Surface) {
	entity.X = (surface.W - entity.W) / 2
}

// CenterEntity centers an entity within a given width
func CenterEntity(entity *Entity, width int32) {
	entity.X = (width - entity.W) / 2
}

// RightJustifyEntityInParent right justifies an entity within its parent
func RightJustifyEntityInParent(entity, parent *Entity) {
	entity.X = parent.W - entity.W
}

// RightJustifyEntityInSurface right justifies an entity within a surface
func RightJustifyEntityInSurface(entity *Entity, surface *sdl.Surface) {
	entity.X = surface.W - entity.W
}

// RightJustifyEntity right justifies an entity within a given width
func RightJustifyEntity(entity *Entity, width int32) {
	entity.X = width - entity.W
}

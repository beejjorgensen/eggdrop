package menu

import (
	"fmt"

	"github.com/beejjorgensen/eggdrop/assetmanager"
	"github.com/beejjorgensen/eggdrop/scenegraph"
	"github.com/beejjorgensen/eggdrop/util"
	"github.com/veandco/go-sdl2/sdl"
)

// Menu holds information about on-screen menus
type Menu struct {
	items             []Item
	selected, clicked int
	spacing           int32 // px
	justification     int
	RootEntity        *scenegraph.Entity
}

// Item describes an individual Menu line item
type Item struct {
	AssetFontID    int
	Text           string
	Color, HiColor sdl.Color
}

// Justification constants for Menu
const (
	MenuJustifyLeft = iota
	MenuJustifyCenter
	MenuJustifyRight
)

// New constructs the menu
func New(am *assetmanager.AssetManager, id int, items []Item, spacing int32, justification int) *Menu {
	// TODO: do these even need to be registered with any asset manager? just use util.RenderText?
	menu := &Menu{items: items, spacing: spacing, justification: justification}

	menuEntity := scenegraph.NewEntity(nil)

	maxH := int32(0)
	maxW := int32(0)

	for i, item := range menu.items {
		var err error
		var surface, surfaceHi *sdl.Surface
		var entity, entityHi *scenegraph.Entity

		if surface, err = am.RenderText(id, item.AssetFontID, item.Text, item.Color); err != nil {
			panic(fmt.Sprintf("Menu render font: %v", err))
		}

		id++

		if surfaceHi, err = am.RenderText(id, item.AssetFontID, item.Text, item.HiColor); err != nil {
			panic(fmt.Sprintf("Intro render font: %v", err))
		}

		id++

		entity = scenegraph.NewEntity(surface)
		entity.Visible = i > 0
		entityHi = scenegraph.NewEntity(surfaceHi)
		entityHi.Visible = i == 0

		menuEntity.AddChild(entity, entityHi)

		if entity.W > maxW {
			maxW = entity.W
		}

		entity.Y = maxH
		entityHi.Y = maxH

		maxH += menu.spacing
	}

	menuEntity.W = maxW
	menuEntity.H = maxH

	// position everything now that we have sizes known
	for i := range menu.items {

		entity := menuEntity.GetChild(i * 2)
		entityHi := menuEntity.GetChild(i*2 + 1)

		switch menu.justification {
		case MenuJustifyLeft:
			entity.X = 0
			entityHi.X = 0
		case MenuJustifyCenter:
			util.CenterEntityInParent(entity, menuEntity)
			util.CenterEntityInParent(entityHi, menuEntity)
		case MenuJustifyRight:
			util.RightJustifyEntityInParent(entity, menuEntity)
			util.RightJustifyEntityInParent(entityHi, menuEntity)
		}
	}

	menu.RootEntity = menuEntity

	return menu
}

// updateVisibility sets the visibility flags on the appropriate elements
func (m *Menu) updateVisibility() {
	root := m.RootEntity

	for i := range m.items {
		eIndex := i * 2

		if i == m.selected {
			root.GetChild(eIndex).Visible = false    // unselected
			root.GetChild(eIndex + 1).Visible = true // selected
		} else {
			root.GetChild(eIndex).Visible = true
			root.GetChild(eIndex + 1).Visible = false
		}

	}
}

// SetSelected sets the selected item in the menu
func (m *Menu) SetSelected(i int) {
	m.selected = i
}

// GetSelected returns the selected item in the menu
func (m *Menu) GetSelected() int {
	return m.selected
}

// SelectNext selects the next item in the menu
func (m *Menu) SelectNext() {
	m.selected++

	if m.selected >= len(m.items) {
		m.selected = 0
	}

	m.updateVisibility()
}

// SelectPrev selects the previous item in the menu
func (m *Menu) SelectPrev() {
	m.selected--

	if m.selected < 0 {
		m.selected = len(m.items) - 1
	}

	m.updateVisibility()
}

// SelectByMouseY selects based on the Y position given
func (m *Menu) SelectByMouseY(y int32) {
	root := m.RootEntity

	y += root.WorldToEntity.Y

	for i := range m.items {
		eIndex := i * 2

		child := root.GetChild(eIndex)

		if y >= child.Y && y <= child.Y+child.H {
			m.selected = i
			m.updateVisibility()
			break
		}
	}
}

// SelectByMouseClickY selects based on the Y position given
func (m *Menu) SelectByMouseClickY(y int32) {
	root := m.RootEntity

	y += root.WorldToEntity.Y

	m.clicked = -1

	for i := range m.items {
		eIndex := i * 2

		child := root.GetChild(eIndex)

		if y >= child.Y && y <= child.Y+child.H {
			m.clicked = i
		}
	}
}

// GetClicked gets the most recently clicked menu entry
func (m *Menu) GetClicked() int {
	return m.clicked
}

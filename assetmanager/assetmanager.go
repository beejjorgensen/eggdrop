// Package assetmanager loads assets (images and fonts) so they can be referred
// to later by ID.
package assetmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/beejjorgensen/eggdrop/util"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var assetDir string

// AssetManager holds surfaces and other asset information
type AssetManager struct {
	Surfaces map[string]*sdl.Surface
	Fonts    map[string]*ttf.Font

	outerSurface *sdl.Surface
}

// AssetPath searches and stores the asset directory. See searchdirs.go for
// paths.
func AssetPath(assetFile string) string {
	if assetDir == "" {
		for _, dir := range searchDirs {
			_, err := os.Stat(dir)
			if err == nil {
				// No error? This must be it. We presume. We hope. We pray.
				assetDir = dir
				break
			}
		}
	}

	if assetDir == "" {
		// Still don't have it? Surrender!
		panic("Assets not found")
	}

	return fmt.Sprintf("%s%c%s", assetDir, os.PathSeparator, assetFile)
}

// New creates and initializes a new AssetManager
func New() *AssetManager {
	return &AssetManager{
		Surfaces: make(map[string]*sdl.Surface),
		Fonts:    make(map[string]*ttf.Font),
	}
}

// SetOuterSurface takes a reference to the outermost surface so it can be used
// later for width and height values from the LoadJSON method. If it's not
// needed in the JSON, this function doesn't need to be called.
func (am *AssetManager) SetOuterSurface(surface *sdl.Surface) {
	am.outerSurface = surface
}

// LoadSurface loads and tracks a new image
func (am *AssetManager) LoadSurface(key string, fileName string) (surface *sdl.Surface, err error) {
	surface, err = img.Load(fileName)
	am.Surfaces[key] = surface

	return
}

// AddSurface tracks an existing surface
func (am *AssetManager) AddSurface(key string, surface *sdl.Surface) {
	am.Surfaces[key] = surface
}

// LoadFont loads and tracks a Font
func (am *AssetManager) LoadFont(key string, fileName string, size int) (err error) {
	am.Fonts[key], err = ttf.OpenFont(fileName, size)

	return
}

// RenderText is a helper function to generate and track a surface with some text on it
func (am *AssetManager) RenderText(surfaceKey, fontKey string, text string, color sdl.Color) (*sdl.Surface, error) {
	surface, err := util.RenderText(am.Fonts[fontKey], text, color)

	if err == nil {
		am.AddSurface(surfaceKey, surface)
	}

	return surface, err
}

// loadJSONPanic has pretty panic events for LoadJSON
func loadJSONPanic(fileName, id, message string) {
	if id == "" {
		panic(fmt.Sprintf("AssetManager: LoadJSON(\"%s\"): %s", fileName, message))
	}
	panic(fmt.Sprintf("AssetManager: LoadJSON(\"%s\"): %s: %s", fileName, id, message))
}

// loadJSONFonts handles loading a Fonts array from the JSON
func (am *AssetManager) loadJSONFonts(jsonFile string, data interface{}) {
	fontArray := data.([]interface{})

	for _, v := range fontArray {

		fontInfo := v.(map[string]interface{})

		id, ok := fontInfo["Id"].(string)
		if !ok {
			loadJSONPanic(jsonFile, "", "Missing ID")
		}
		font, ok := fontInfo["Font"].(string)
		if !ok {
			loadJSONPanic(jsonFile, id, "Missing Font")
		}
		size, ok := fontInfo["Size"].(float64)
		if !ok {
			loadJSONPanic(jsonFile, id, "Missing Size")
		}

		if err := am.LoadFont(id, AssetPath(font), int(size)); err != nil {
			loadJSONPanic(jsonFile, id, fmt.Sprintf("%v", err))
		}
	}
}

// loadJSONImages handles loading an Images array from the JSON
func (am *AssetManager) loadJSONImages(jsonFile string, data interface{}) {
	imageArray := data.([]interface{})

	for _, v := range imageArray {
		imageInfo := v.(map[string]interface{})

		id, ok := imageInfo["Id"].(string)
		if !ok {
			loadJSONPanic(jsonFile, "", "Image missing ID")
		}

		image, imageOk := imageInfo["Image"].(string)
		src, srcOk := imageInfo["Src"].(string)

		if imageOk {
			// Load a normal image
			_, err := am.LoadSurface(id, AssetPath(image))
			if err != nil {
				loadJSONPanic(jsonFile, id, "Error loading image")
			}

		} else if srcOk {
			// Generate an image from existing with some transformation
			if _, ok := am.Surfaces[src]; !ok {
				loadJSONPanic(jsonFile, id, "Unknown Src")
			}

			transform, _ := imageInfo["Transform"]
			switch transform {
			//case "COPY":
			case "FLIP_V":
				newSurface, err := util.SurfaceFlipV(am.Surfaces[src])
				if err != nil {
					loadJSONPanic(jsonFile, id, "Error on vertical flip")
				}
				am.AddSurface(id, newSurface)
			case "FLIP_H":
				newSurface, err := util.SurfaceFlipH(am.Surfaces[src])
				if err != nil {
					loadJSONPanic(jsonFile, id, "Error on horizontal flip")
				}
				am.AddSurface(id, newSurface)
			default:
				loadJSONPanic(jsonFile, id, "Unrecognized Transform")
			}

		} else {
			loadJSONPanic(jsonFile, id, "Must specify Image or Src")
		}
	}

}

// sdlColorFromInterfaceArray is a helper function for parsing RGBA arrays from
// JSON.
func sdlColorFromInterfaceArray(rgba []interface{}) (*sdl.Color, error) {
	r, rOk := rgba[0].(float64)
	g, gOk := rgba[1].(float64)
	b, bOk := rgba[2].(float64)
	a, aOk := rgba[3].(float64)

	if !rOk || !gOk || !bOk || !aOk {
		return nil, errors.New("RGBA parse failure")
	}

	return &sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}, nil
}

// renderJSONText renders text from the JSON file into surfaces
func (am *AssetManager) renderJSONText(jsonFile string, data interface{}) {
	var err error

	textArray := data.([]interface{})
	var color *sdl.Color

	for _, v := range textArray {
		textInfo := v.(map[string]interface{})

		id, ok := textInfo["Id"].(string)
		if !ok {
			loadJSONPanic(jsonFile, "", "Text missing Id")
		}
		font, ok := textInfo["Font"].(string)
		if !ok {
			loadJSONPanic(jsonFile, id, "Missing Font")
		}
		text, ok := textInfo["Text"].(string)
		if !ok {
			loadJSONPanic(jsonFile, id, "Missing Text")
		}
		rgba, ok := textInfo["Rgba"].([]interface{})
		if len(rgba) == 4 {
			color, err = sdlColorFromInterfaceArray(rgba)
		}
		if !ok || len(rgba) != 4 || err != nil {
			loadJSONPanic(jsonFile, id, "Rgba needs to be in form [R,G,B,A], 0-255 for each element")
		}

		if _, err = am.RenderText(id, font, text, *color); err != nil {
			panic(fmt.Sprintf("Intro render font: %v", err))
		}
	}
}

// getWindowInfoFromString takes a string from the JSON that requests window
// size info and returns the proper size
func (am *AssetManager) getWindowInfoFromString(s string) (int32, error) {
	if am.outerSurface == nil {
		return 0, errors.New("outerSurface not specified")
	}
	switch s {
	case "WINDOW_WIDTH":
		return int32(am.outerSurface.W), nil
	case "WINDOW_HEIGHT":
		return int32(am.outerSurface.H), nil
	default:
		return 0, fmt.Errorf("unrecognized parameter: %s", s)
	}
}

// parseRectDimension parses a W or H dimension out of a Rect Record
func (am *AssetManager) parseRectDimension(jsonFile, id, dim string, info map[string]interface{}) int32 {
	var rv int32
	var err error

	switch val := info[dim].(type) { // "W" or "H"
	case float64:
		rv = int32(val)
	case string:
		rv, err = am.getWindowInfoFromString(val)
		if err != nil {
			loadJSONPanic(jsonFile, id, fmt.Sprintf("%s: %v", val, err))
		}
	default:
		loadJSONPanic(jsonFile, id, fmt.Sprintf("%t: unknown parameter type", val))
	}

	return rv
}

// renderJSONRects renders rectangles specified in the JSON assets file
func (am *AssetManager) renderJSONRects(jsonFile string, data interface{}) {
	var err error

	rectsArray := data.([]interface{})
	var color *sdl.Color

	for _, v := range rectsArray {
		rectInfo := v.(map[string]interface{})

		id, ok := rectInfo["Id"].(string)
		if !ok {
			loadJSONPanic(jsonFile, "", "Rect missing Id")
		}

		rgba, ok := rectInfo["Rgba"].([]interface{})
		if len(rgba) == 4 {
			color, err = sdlColorFromInterfaceArray(rgba)
		}
		if !ok || len(rgba) != 4 || err != nil {
			loadJSONPanic(jsonFile, id, "Rgba needs to be in form [R,G,B,A], 0-255 for each element")
		}

		var width, height int32

		width = am.parseRectDimension(jsonFile, id, "W", rectInfo)
		height = am.parseRectDimension(jsonFile, id, "H", rectInfo)

		surface, err := util.MakeFillSurfaceAlpha(width, height, color.R, color.G, color.B, color.A)
		if err != nil {
			loadJSONPanic(jsonFile, id, fmt.Sprintf("Error making rect: %v", err))
		}

		am.AddSurface(id, surface)
	}
}

// LoadJSON reads a JSON file and loads all the assets described therein
func (am *AssetManager) LoadJSON(jsonFile string) error {
	var err error

	jsonStr, err := ioutil.ReadFile(AssetPath(jsonFile))

	var assetData map[string]interface{} // outer JSON is an object

	err = json.Unmarshal(jsonStr, &assetData)
	if err != nil {
		return err
	}

	// Process fonts first since Text nodes can refer to them
	fonts, ok := assetData["Fonts"]
	if ok {
		am.loadJSONFonts(jsonFile, fonts)
	}

	for k, v := range assetData {
		switch k {
		case "Fonts":
			// Do nothing--we processed them above

		case "Images":
			am.loadJSONImages(jsonFile, v)

		case "Text":
			am.renderJSONText(jsonFile, v)

		case "Rects":
			am.renderJSONRects(jsonFile, v)
		}

	}

	return nil
}

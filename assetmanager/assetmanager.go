// Package assetmanager loads assets (images and fonts) so they can be referred
// to later by ID.
package assetmanager

import (
	"encoding/json"
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

// LoadJSON reads a JSON file and loads all the assets described therein
func (am *AssetManager) LoadJSON(jsonFile string) error {
	var err error

	jsonStr, err := ioutil.ReadFile(AssetPath(jsonFile))

	var assetData map[string]interface{} // outer JSON is an object

	err = json.Unmarshal(jsonStr, &assetData)
	if err != nil {
		return err
	}

	for k, v := range assetData {
		switch k {
		case "Fonts":
			am.loadJSONFonts(jsonFile, v)
		}
	}

	return nil
}

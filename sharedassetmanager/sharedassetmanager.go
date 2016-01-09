// Package sharedassetmanager declares a global AssetManager that could be used
// by multiple packages. Right now the main game modes have their own
// AssetManagers, but that precludes them from sharing assets.
package sharedassetmanager

import "beej.us/eggdrop/assetmanager"

// GAssetManager is a global asset manager for shared resources
var GAssetManager = assetmanager.New()

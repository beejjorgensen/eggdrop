package assetmanager

// This is out here for easy patching by install scripts
var searchDirs = []string{
	// --ADD NEW SEARCH PATH HERE--
	"/usr/lib/eggdrop/assets",
	"/usr/share/eggdrop/assets",
	"/usr/local/lib/eggdrop/assets",
	"/usr/local/share/eggdrop/assets",
	"./assets",
}

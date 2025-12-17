package curseforge

import "fmt"

// Loader type constants
const (
	BukkitLoaderID   = 1     // Bukkit/Spigot/Paper plugins
	FabricLoaderID   = 4
	ForgeLoaderID    = 1
	NeoForgeLoaderID = 6
	QuiltLoaderID    = 5
)

// Game version type IDs for different Minecraft versions
const (
	PluginVersionType  = 1     // Bukkit/Plugin versions (all versions)
	ModVersionType_119 = 73407 // Mod loader versions for 1.19.x
	ModVersionType_120 = 75125 // Mod loader versions for 1.20.x
	ModVersionType_121 = 77784 // Mod loader versions for 1.21.x
)

// pluginVersionToID maps Minecraft version strings to CurseForge game version IDs
// These IDs are for gameVersionTypeID = 1 (Bukkit/Plugin versions)
var pluginVersionToID = map[string]int{
	// 1.19.x versions
	"1.19":   9190,
	"1.19.1": 9261,
	"1.19.2": 9560,
	"1.19.3": 9561,
	"1.19.4": 9973,

	// 1.20.x versions
	"1.20":   9974,
	"1.20.1": 9994,
	"1.20.2": 10326,
	"1.20.3": 10741,
	"1.20.4": 10742,
	"1.20.5": 11306,
	"1.20.6": 11307,

	// 1.21.x versions
	"1.21":    11515,
	"1.21.1":  12735,
	"1.21.2":  12736,
	"1.21.3":  12737,
	"1.21.4":  12738,
	"1.21.5":  12988,
	"1.21.6":  13473,
	"1.21.7":  13574,
	"1.21.8":  13683,
	"1.21.9":  13933,
	"1.21.10": 13966,
	"1.21.11": 14417,
}

// modVersionToID maps Minecraft version strings to CurseForge game version IDs
// These IDs are for mod loader versions (gameVersionTypeID varies by MC version)
var modVersionToID = map[string]int{
	// 1.19.x versions (gameVersionTypeID = 73407)
	"1.19":   9186,
	"1.19.1": 9259,
	"1.19.2": 9366,
	"1.19.3": 9550,
	"1.19.4": 9776,

	// 1.20.x versions (gameVersionTypeID = 75125)
	"1.20":   9971,
	"1.20.1": 9990,
	"1.20.2": 10236,
	"1.20.3": 10395,
	"1.20.4": 10407,
	"1.20.5": 11163,
	"1.20.6": 11198,

	// 1.21.x versions (gameVersionTypeID = 77784)
	"1.21":    11457,
	"1.21.1":  11779,
	"1.21.2":  12079,
	"1.21.3":  12084,
	"1.21.4":  12281,
	"1.21.5":  12934,
	"1.21.6":  13422,
	"1.21.7":  13506,
	"1.21.8":  13620,
	"1.21.9":  13927,
	"1.21.10": 13964,
	"1.21.11": 14406,
}

// ConvertVersionString converts a Minecraft version string to CurseForge ID
// projectType should be "plugin" or "mod"
func ConvertVersionString(version string, projectType string) (int, bool) {
	if projectType == "mod" {
		if id, exists := modVersionToID[version]; exists {
			return id, true
		}
	} else {
		// Default to plugin
		if id, exists := pluginVersionToID[version]; exists {
			return id, true
		}
	}
	return 0, false
}

// GetLoaderID returns the CurseForge loader ID for a given project type and loader
func GetLoaderID(projectType string, loader string) (int, error) {
	if projectType == "plugin" {
		return BukkitLoaderID, nil
	}

	switch loader {
	case "fabric":
		return FabricLoaderID, nil
	case "forge":
		return ForgeLoaderID, nil
	case "neoforge":
		return NeoForgeLoaderID, nil
	case "quilt":
		return QuiltLoaderID, nil
	default:
		return 0, fmt.Errorf("unknown loader type: %s", loader)
	}
}
package curseforge

import "fmt"

// Loader IDs for CurseForge
const (
	BukkitLoaderID = 5  // Bukkit/Spigot/Paper plugins
	FabricLoaderID = 4  // Fabric mods
	ForgeLoaderID  = 1  // Forge mods
	NeoForgeLoaderID = 6 // NeoForge mods
	QuiltLoaderID  = 5  // Quilt mods (note: may conflict with Bukkit, verify if needed)
)

// minecraftVersionToID maps Minecraft version strings to CurseForge game version IDs
var minecraftVersionToID = map[string]int{
	"1.21":     10407,
	"1.21.1":   10785,
	"1.21.2":   11092,
	"1.21.3":   11213,
	"1.21.4":   11596,
	"1.20.6":   10235,
	"1.20.5":   10169,
	"1.20.4":   9971,
	"1.20.3":   9883,
	"1.20.2":   9856,
	"1.20.1":   9990,
	"1.20":     9885,
	"1.19.4":   9776,
	"1.19.3":   9550,
	"1.19.2":   9366,
	"1.19.1":   9259,
	"1.19":     9186,
	"1.18.2":   9008,
	"1.18.1":   8857,
	"1.18":     8830,
	"1.17.1":   8516,
	"1.17":     8203,
	"1.16.5":   7915,
	"1.16.4":   7890,
	"1.16.3":   7667,
	"1.16.2":   7498,
	"1.16.1":   7498,
	"1.16":     7469,
}

// ConvertVersionString converts a Minecraft version string to CurseForge ID
func ConvertVersionString(version string) (int, bool) {
	if id, exists := minecraftVersionToID[version]; exists {
		return id, true
	}
	return 0, false
}

// GetLoaderID returns the appropriate loader ID based on the project type and loader
func GetLoaderID(projectType, loader string) (int, error) {
	if projectType == "plugin" || projectType == "" { // Default to plugin for backward compatibility
		return BukkitLoaderID, nil
	}
	
	if projectType == "mod" {
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
			return 0, fmt.Errorf("unknown mod loader: %s (supported: fabric, forge, neoforge, quilt)", loader)
		}
	}
	
	return 0, fmt.Errorf("unknown project type: %s (supported: plugin, mod)", projectType)
}
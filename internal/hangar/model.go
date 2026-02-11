package hangar

type VersionUploadReq struct {
	Version              string                          `json:"version"`
	PluginDependencies   map[Platform][]PluginDependency `json:"pluginDependencies"`
	PlatformDependencies map[Platform][]string           `json:"platformDependencies"`
	Description          string                          `json:"description"`
	Files                []MultipartFileOrURL            `json:"files"`
	Channel              string                          `json:"channel"`
}

type Platform string

const (
	PlatformPaper     Platform = "PAPER"
	PlatformWaterfall Platform = "WATERFALL"
	PlatformVelocity  Platform = "VELOCITY"
)

type PluginDependency struct {
	Name        string  `json:"name"`
	Required    bool    `json:"required"`
	ExternalURL *string `json:"externalUrl,omitempty"` // nullable
}

type MultipartFileOrURL struct {
	Platforms   []Platform `json:"platforms"`
	ExternalURL *string    `json:"externalUrl,omitempty"` // nullable
}

type AuthenticateResp struct {
	Token string `json:"token"`
}

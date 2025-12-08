package orbis

type CreateVersionReq struct {
	VersionNumber      string   `json:"versionNumber"`
	Name               string   `json:"name"`
	Channel            string   `json:"channel"`
	CompatibleVersions []string `json:"compatibleVersions"`
	Changelog          string   `json:"changelog"`
}

type Version struct {
	ID string `json:"id"`
}

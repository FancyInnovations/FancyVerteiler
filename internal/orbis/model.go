package orbis

type CreateVersionReq struct {
	VersionNumber              string   `json:"versionNumber"`
	Name                       string   `json:"name"`
	Channel                    string   `json:"channel"`
	CompatibleHytaleVersionIds []string `json:"compatibleHytaleVersionIds"`
}

type UpdateChangelogReq struct {
	Changelog string `json:"changelog"`
}

type SetPrimaryFileReq struct {
	FileID string `json:"fileId"`
}

type Version struct {
	ID string `json:"id"`
}

type VersionFile struct {
	ID string `json:"id"`
}

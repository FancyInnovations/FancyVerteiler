package orbis

type CreateVersionReq struct {
	VersionNumber              string   `json:"versionNumber"`
	Channel                    string   `json:"channel"`
	CompatibleHytaleVersionIds []string `json:"compatibleHytaleVersionIds"`
}

type UpdateChangelogReq struct {
	Changelog string `json:"changelog"`
}

type SetPrimaryFileReq struct {
	FileID string `json:"fileId"`
}

type VersionResp struct {
	Version Version `json:"version"`
}
type Version struct {
	ID string `json:"id"`
}

type UploadFileResp struct {
	File VersionFile `json:"file"`
}

type VersionFile struct {
	ID string `json:"id"`
}

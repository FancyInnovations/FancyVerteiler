package curseforge

type CreateVersionReq struct {
	Changelog     string                 `json:"changelog"`
	ChangelogType string                 `json:"changelogType"`
	DisplayName   string                 `json:"displayName"`
	GameVersions  []int                  `json:"gameVersions"`
	ReleaseType   string                 `json:"releaseType"`
	Relations     CreateVersionRelations `json:"relations"`
}

type CreateVersionRelations struct {
	Projects []ProjectRelation `json:"projects"`
}

type ProjectRelation struct {
	Slug string `json:"slug"`
	Type string `json:"type"`
}
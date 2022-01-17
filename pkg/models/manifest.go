package models

type Manifest struct {
	BuildNumber string
	Repo        struct {
		BranchName string
		CommitSha  string
	}
	Artefacts map[string]string
	MetaData  ManifestMetaData
}

type ManifestMetaData struct {
	BldrVersion     string
	ManifestVersion string
}

package contexts

type BuildContext struct {
	BuildNumber      string
	GitContext       *GitContext
	ArtefactManifest *ArtefactManifest
}

func CreateBuildContext() *BuildContext {
	return &BuildContext{
		GitContext: &GitContext{},
		ArtefactManifest: &ArtefactManifest{
			Artefacts: make(map[string]string),
		},
	}
}

type GitContext struct {
	FullCommitSha       string
	BranchName          string
	MainBranchForkPoint string
}

type ArtefactManifest struct {
	Artefacts map[string]string
}

func (am *ArtefactManifest) AddArtefact(targetName string, artefactPath string) {
	am.Artefacts[targetName] = artefactPath
}

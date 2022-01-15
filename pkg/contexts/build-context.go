package contexts

type BuildContext struct {
	BuildNumber      string
	ArtefactManifest *ArtefactManifest
	GitContext       *GitContext
}

func CreateBuildContext() *BuildContext {
	return &BuildContext{
		ArtefactManifest: &ArtefactManifest{
			Artefacts: make(map[string]string),
		},
		GitContext: &GitContext{},
	}
}

type ArtefactManifest struct {
	Artefacts map[string]string
}

type GitContext struct {
	FullCommitSha            string
	ShortCommitSha           string
	BranchName               string
	MainBranchForkPoint      string
	MainBranchForkPointShort string
	ChangesSinceMainBranch   []string
}

func (am *ArtefactManifest) AddArtefact(targetName string, artefactPath string) {
	am.Artefacts[targetName] = artefactPath
}

func (g *GitContext) CanDetectChanges() bool {
	return len(g.ChangesSinceMainBranch) > 0
}

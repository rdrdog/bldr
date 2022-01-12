package contexts

type BuildContext struct {
	BuildNumber      string
	ArtefactManifest *ArtefactManifest
	DockerContext    *DockerContext
	GitContext       *GitContext
	PathContext      *PathContext
}

func CreateBuildContext( /*TODO - set paths */ ) *BuildContext {
	return &BuildContext{
		ArtefactManifest: &ArtefactManifest{
			Artefacts: make(map[string]string),
		},
		DockerContext: &DockerContext{
			UseBuildKit:                     true,
			IncludeTimeInTags:               false,
			UseRemoteContainerRegistryCache: false,
		},
		GitContext:  &GitContext{},
		PathContext: &PathContext{},
	}
}

type ArtefactManifest struct {
	Artefacts map[string]string
}

type DockerContext struct {
	UseBuildKit                     bool
	Registry                        string
	IncludeTimeInTags               bool
	UseRemoteContainerRegistryCache bool
}

type GitContext struct {
	FullCommitSha       string
	ShortCommitSha      string
	BranchName          string
	MainBranchForkPoint string
}

type PathContext struct {
	RepoRootDirectory string
}

func (am *ArtefactManifest) AddArtefact(targetName string, artefactPath string) {
	am.Artefacts[targetName] = artefactPath
}

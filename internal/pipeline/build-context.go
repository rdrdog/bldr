package pipeline

type BuildContext struct {
	BuildNumber string
	GitContext  *GitContext
}

type GitContext struct {
	FullCommitSha       string
	BranchName          string
	MainBranchForkPoint string
}

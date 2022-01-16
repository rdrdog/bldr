package lib

type GitState struct {
	CommitSha              string
	BranchName             string
	MainBranchForkPoint    string
	ChangesSinceMainBranch []string
}

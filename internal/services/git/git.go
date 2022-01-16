package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/rdrdog/bldr/internal/services/process"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/sirupsen/logrus"
)

const EnvVarGitForkPoint = "GIT_MAIN_BRANCH_FORK_COMMIT"
const EnvVarGitBranchName = "GIT_BRANCH_NAME"

type Git struct {
	logger            *logrus.Logger
	repoRootDirectory string
	mainBranchName    string
}

func New(logger *logrus.Logger, mainBranchName string, repoRootDirectory string) *Git {
	return &Git{
		logger:            logger,
		mainBranchName:    mainBranchName,
		repoRootDirectory: repoRootDirectory,
	}
}

func (g *Git) getRepoCommitSha() (string, error) {

	commitSha, _, err := process.
		New("git", g.repoRootDirectory, g.logger).
		WithArgs("rev-parse HEAD").
		WithSuppressedOutput().
		Run()

	return commitSha, err
}

func (g *Git) getMainBranchForkPoint() (string, error) {

	injectedForkPoint := os.Getenv(EnvVarGitForkPoint)
	if len(injectedForkPoint) > 0 {
		g.logger.Infof("Env var override found for %s", EnvVarGitForkPoint)
		return injectedForkPoint, nil
	}

	commitSha, _, err := process.
		New("git", g.repoRootDirectory, g.logger).
		WithArgs("merge-base --octopus remotes/origin/" + g.mainBranchName + " HEAD").
		WithSuppressedOutput().
		Run()

	return commitSha, err
}

func (g *Git) getChangesBetweenCommits(originCommitSha string, currentCommitSha string) ([]string, error) {

	var args string
	if true /*g.context.IncludeLocalChangesInDiff*/ {
		args = fmt.Sprintf("--no-pager diff --name-only %s:./", originCommitSha)
	} else {
		args = fmt.Sprintf("--no-pager diff --name-only %s..%s", originCommitSha, currentCommitSha)
	}

	diffOutput, _, procError := process.
		New("git", g.repoRootDirectory, g.logger).
		WithArgs(args).
		WithSuppressedOutput().
		Run()

	return strings.Split(diffOutput, "\n"), procError
}

func (g *Git) getBranchName() (string, error) {

	injectedBranchName := os.Getenv(EnvVarGitBranchName)
	if len(injectedBranchName) > 0 {
		g.logger.Infof("Env var override found for %s", EnvVarGitBranchName)
		return injectedBranchName, nil
	}

	branchName, _, err := process.
		New("git", g.repoRootDirectory, g.logger).
		WithArgs("rev-parse --abbrev-ref HEAD").
		WithSuppressedOutput().
		Run()

	return branchName, err
}

func (g *Git) LoadRepoInformation() *lib.GitState {
	result := &lib.GitState{}
	var err error
	result.CommitSha, err = g.getRepoCommitSha()
	if err != nil {
		g.logger.WithField("error", err).Warn("Error getting commit sha")
	}

	result.BranchName, err = g.getBranchName()
	if err != nil {
		g.logger.WithField("error", err).Warn("Error getting branch name")
	}

	result.MainBranchForkPoint, err = g.getMainBranchForkPoint()
	if err != nil {
		g.logger.WithField("error", err).Warn("Error getting main branch fork point")
	}

	if len(result.CommitSha) > 0 && len(result.MainBranchForkPoint) > 0 {
		result.ChangesSinceMainBranch, err = g.getChangesBetweenCommits(result.MainBranchForkPoint, result.CommitSha)
		if err != nil {
			g.logger.WithField("error", err).Warn("Error getting diff list")
		} else {
			g.logger.Debugf("Diff for branch: %s", result.ChangesSinceMainBranch)
		}
	}

	g.logger.
		WithField("commitSha", result.CommitSha).
		WithField("branchName", result.BranchName).
		WithField("mainBranchForkPoint", result.MainBranchForkPoint).
		Info("Loaded repo information")

	return result
}

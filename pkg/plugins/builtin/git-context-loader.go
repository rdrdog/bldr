package builtin

import (
	"github.com/Redgwell/bldr/pkg/contexts"
	"github.com/Redgwell/bldr/pkg/lib/git"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type GitContextLoader struct {
	logger *logrus.Logger
	Name   string
}

func (p *GitContextLoader) SetConfig(logger *logrus.Logger, targetName string, pluginConfig map[string]interface{}) error {
	p.logger = logger
	p.Name = targetName
	return mapstructure.Decode(pluginConfig, p)
}

func (p *GitContextLoader) Execute(contextProvider *contexts.ContextProvider) error {
	bc := contextProvider.BuildContext
	p.logger.Infof("Loading git context for path %s", bc.PathContext.RepoRootDirectory)

	git := git.New(p.logger, bc.GitContext.MainBranchName, bc.PathContext.RepoRootDirectory)
	git.LoadRepoInformation()

	bc.GitContext.BranchName = git.BranchName
	bc.GitContext.FullCommitSha = git.CommitSha
	bc.GitContext.ShortCommitSha = bc.GitContext.FullCommitSha[:7]
	bc.GitContext.MainBranchForkPoint = git.MainBranchForkPoint
	bc.GitContext.ChangesSinceMainBranch = git.ChangesSinceMainBranch

	return nil
}

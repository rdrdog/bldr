package builtin

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/lib/git"
	"github.com/sirupsen/logrus"
)

type GitContextLoader struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	Name          string
}

func (p *GitContextLoader) SetConfig(logger *logrus.Logger, targetName string, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	p.Name = targetName
	return nil
}

func (p *GitContextLoader) Execute(contextProvider *contexts.ContextProvider) error {
	bc := contextProvider.BuildContext
	p.logger.Infof("loading git context for path %s", bc.PathContext.RepoRootDirectory)

	git := git.New(p.logger, p.configuration.Git.MainBranchName, bc.PathContext.RepoRootDirectory)
	git.LoadRepoInformation()

	bc.GitContext.BranchName = git.BranchName
	bc.GitContext.FullCommitSha = git.CommitSha
	bc.GitContext.ShortCommitSha = git.CommitSha[:7]
	bc.GitContext.MainBranchForkPointShort = git.MainBranchForkPoint[:7]
	bc.GitContext.ChangesSinceMainBranch = git.ChangesSinceMainBranch

	return nil
}

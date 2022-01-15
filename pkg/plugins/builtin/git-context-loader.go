package builtin

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib/git"
	"github.com/sirupsen/logrus"
)

type GitContextLoader struct {
	configuration *config.Configuration
	logger        *logrus.Logger
}

func (p *GitContextLoader) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	return nil
}

func (p *GitContextLoader) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider) error {
	bc := contextProvider.GetBuildContext()
	p.logger.Infof("loading git context for path %s", p.configuration.Paths.RepoRootDirectory)

	git := git.New(p.logger, p.configuration.Git.MainBranchName, p.configuration.Paths.RepoRootDirectory)
	git.LoadRepoInformation()

	bc.GitContext.BranchName = git.BranchName
	bc.GitContext.FullCommitSha = git.CommitSha
	bc.GitContext.ShortCommitSha = git.CommitSha[:7]
	bc.GitContext.MainBranchForkPointShort = git.MainBranchForkPoint[:7]
	bc.GitContext.ChangesSinceMainBranch = git.ChangesSinceMainBranch

	return nil
}

package builtin

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib"
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

func (p *GitContextLoader) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider, libProvider lib.LibProvider) error {
	bc := contextProvider.GetBuildContext()

	gitState := libProvider.GetGitLib().LoadRepoInformation()

	bc.GitContext.BranchName = gitState.BranchName
	bc.GitContext.FullCommitSha = gitState.CommitSha
	bc.GitContext.ShortCommitSha = gitState.CommitSha[:7]
	bc.GitContext.MainBranchForkPointShort = gitState.MainBranchForkPoint[:7]
	bc.GitContext.ChangesSinceMainBranch = gitState.ChangesSinceMainBranch

	return nil
}

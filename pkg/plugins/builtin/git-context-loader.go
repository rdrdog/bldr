package builtin

import (
	"github.com/Redgwell/bldr/pkg/contexts"
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

	// TODO - load:
	// bc.GitContext.BranchName
	// bc.GitContext.FullCommitSha
	// bc.GitContext.ShortCommitSha = bc.GitContext.FullCommitSha[:7]
	// bc.GitContext.MainBranchForkPoint

	return nil
}

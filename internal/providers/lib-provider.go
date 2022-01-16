package providers

import (
	"github.com/rdrdog/bldr/internal/services/docker"
	"github.com/rdrdog/bldr/internal/services/git"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/sirupsen/logrus"
)

type DefaultLibProvider struct {
	logger        *logrus.Logger
	configuration *config.Configuration
}

func NewDefaultLibProvider(logger *logrus.Logger, configuration *config.Configuration) *DefaultLibProvider {
	return &DefaultLibProvider{
		logger:        logger,
		configuration: configuration,
	}
}

func (p *DefaultLibProvider) GetDockerLib() lib.Docker {
	return docker.New(p.configuration, p.logger)
}

func (p *DefaultLibProvider) GetGitLib() lib.Git {
	return git.New(p.logger, p.configuration.Git.MainBranchName, p.configuration.Paths.RepoRootDirectory)
}

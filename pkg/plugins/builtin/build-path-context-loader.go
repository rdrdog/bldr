package builtin

import (
	"os"

	"github.com/Redgwell/bldr/pkg/contexts"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type BuildPathContextLoader struct {
	logger *logrus.Logger
	Name   string
}

func (p *BuildPathContextLoader) SetConfig(logger *logrus.Logger, targetName string, pluginConfig map[string]interface{}) error {
	p.logger = logger
	p.Name = targetName
	return mapstructure.Decode(pluginConfig, p)
}

func (p *BuildPathContextLoader) Execute(contextProvider *contexts.ContextProvider) error {
	pc := contextProvider.BuildContext.PathContext
	p.logger.Info("Loading path context")

	var err error
	pc.RepoRootDirectory, err = os.Getwd()
	if err != nil {
		p.logger.Fatal("could not load working directory in build path context loader")
	}

	p.logger.Debugf("Repo root directory determined to be %s", pc.RepoRootDirectory)
	return err
}

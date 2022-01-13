package builtin

import (
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/sirupsen/logrus"
)

const BuildArtefactDirectoryName = "build-artefacts"

type BuildPathContextLoader struct {
	logger *logrus.Logger
	Name   string
}

func (p *BuildPathContextLoader) SetConfig(logger *logrus.Logger, targetName string, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.logger = logger
	p.Name = targetName
	return mapstructure.Decode(pluginConfig, p)
}

func (p *BuildPathContextLoader) Execute(contextProvider *contexts.ContextProvider) error {
	pc := contextProvider.BuildContext.PathContext
	p.logger.Info("loading path context")

	var err error
	pc.RepoRootDirectory, err = os.Getwd()
	if err != nil {
		p.logger.Fatal("could not load working directory in build path context loader")
		return err
	}

	p.logger.Debugf("repo root directory determined to be %s", pc.RepoRootDirectory)

	pc.BuildArtefactDirectory = path.Join(pc.RepoRootDirectory, BuildArtefactDirectoryName)

	os.RemoveAll(pc.BuildArtefactDirectory)
	err = os.Mkdir(pc.BuildArtefactDirectory, 0755)
	if err != nil {
		p.logger.Fatalf("error creating artefact directory at %s: %v", pc.BuildArtefactDirectory, err)
		return err
	}

	return nil
}

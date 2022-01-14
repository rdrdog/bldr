package builtin

import (
	"os"
	"path"

	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/utils"
	"github.com/sirupsen/logrus"
)

const BuildArtefactDirectoryName = "build-artefacts"
const PipelineConfigFileName = "pipeline-config.yaml"

type BuildPathContextLoader struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	Name          string
}

func (p *BuildPathContextLoader) SetConfig(logger *logrus.Logger, targetName string, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	p.Name = targetName
	return nil
}

func (p *BuildPathContextLoader) Execute(contextProvider *contexts.ContextProvider, extensionsProvider *extensions.ExtensionsProvider) error {
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

	// Copy the pipeline-config.yaml to the artefact directory
	pipelineConfigDst := path.Join(pc.BuildArtefactDirectory, PipelineConfigFileName)
	utils.CopyFile(p.configuration.Pipeline.Path, pipelineConfigDst)

	return nil
}

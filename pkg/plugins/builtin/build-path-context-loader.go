package builtin

import (
	"os"
	"path"

	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/rdrdog/bldr/pkg/utils"
	"github.com/sirupsen/logrus"
)

type BuildPathContextLoader struct {
	configuration *config.Configuration
	logger        *logrus.Logger
}

func (p *BuildPathContextLoader) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	return nil
}

func (p *BuildPathContextLoader) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider, libProvider lib.LibProvider) error {
	p.logger.Info("setting up build paths")

	buildArtefactDir := p.configuration.Paths.BuildArtefactDirectory
	os.RemoveAll(buildArtefactDir)
	err := os.Mkdir(buildArtefactDir, 0755)
	if err != nil {
		p.logger.Fatalf("error creating artefact directory at %s: %v", buildArtefactDir, err)
		return err
	}

	// Copy the pipeline-config.yaml to the artefact directory
	_, pipelineConfigFileName := path.Split(p.configuration.Paths.PipelineConfigFile)
	pipelineConfigDst := path.Join(buildArtefactDir, pipelineConfigFileName)
	utils.CopyFile(p.configuration.Paths.PipelineConfigFile, pipelineConfigDst)

	return nil
}

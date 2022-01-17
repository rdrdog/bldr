package builtin

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/rdrdog/bldr/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

type DeployContextLoader struct {
	configuration   *config.Configuration
	logger          *logrus.Logger
	EnvironmentName string
}

func (p *DeployContextLoader) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	return mapstructure.Decode(pluginConfig, p)
}

func (p *DeployContextLoader) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider, libProvider lib.LibProvider) error {
	p.logger.Info("setting up deployment context")

	dc := contextProvider.GetDeployContext()
	dc.EnvironmentName = p.EnvironmentName

	// Load the manifest
	manifestPath := p.configuration.Paths.DeploymentManifestFile
	exists, err := afero.Exists(config.Appfs, manifestPath)
	if err != nil {
		p.logger.Errorf("error checking if deployment manifest exists at %s: %v", manifestPath, err)
		return err
	}
	if !exists {
		return fmt.Errorf("deployment manifest could not be found at '%s'", manifestPath)
	}

	// Read the manifest config:
	data, err := afero.ReadFile(config.Appfs, manifestPath)
	if err != nil {
		p.logger.Errorf("error reading deployment manifest file from %s: %v", manifestPath, err)
		return err
	}

	// Load the manifest yaml
	manifest := new(models.Manifest)
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		p.logger.Errorf("error loading deployment manifest file yaml from %s: %v", manifestPath, err)
		return err
	}

	dc.Artefacts = manifest.Artefacts
	p.logger.Infof("loading deployment manifest artefacts: %v", dc.Artefacts)

	return nil
}

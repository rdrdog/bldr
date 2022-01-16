package contexts

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

type DeployContext struct {
	config          *config.Configuration
	logger          *logrus.Logger
	EnvironmentName string
	Artefacts       map[string]string
}

func CreateDeployContext(logger *logrus.Logger, config *config.Configuration) *DeployContext {
	// TODO - load artefacts from manifest
	// TODO - load environmentName from *somewhere*
	return &DeployContext{
		config:          config,
		logger:          logger,
		EnvironmentName: "TODO",
		Artefacts:       make(map[string]string), //TODO
	}
}

func (c *DeployContext) GetArtefactByName(name string) string {

	result := c.Artefacts[name]

	if len(result) == 0 {
		c.logger.Fatalf("could not find artefact with name %s in provided manifest", name)
		return ""
	}

	return result
}

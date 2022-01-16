package builtin

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/sirupsen/logrus"
)

type DockerRun struct {
	configuration    *config.Configuration
	logger           *logrus.Logger
	Name             string
	SkipEnvironments []string
	Targets          []DockerRunTargets
}

type DockerRunTargets struct {
	Name    string
	Secrets []interface{}
}

func (p *DockerRun) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	return mapstructure.Decode(pluginConfig, p)
}

func (p *DockerRun) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider) error {
	//dc := contextProvider.DeployContext

	/*
		- Get environment name from deploy context
			- if it's the skip envs, then skip
		- Get manifest from the deploy context
			- get our container tag from the manifest
		- Load secrets using: extensionsProvider.SecretLoader.LoadSecrets()
		- Run docker run...
	*/

	for _, t := range p.Targets {
		secrets, err := extensionsProvider.GetSecretLoader().LoadSecrets(t.Name, t.Secrets)
		if err != nil {
			p.logger.Errorf("error loading secrets for target %s", t.Name)
			return err
		}
		p.logger.Debugf("loaded secrets for target %s", secrets)
	}

	return nil
}

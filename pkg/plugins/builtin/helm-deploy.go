package builtin

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/sirupsen/logrus"
)

type HelmDeploy struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	Name          string
	Targets       []HelmDeployTargets
}

type HelmDeployTargets struct {
	Name    string
	Secrets []interface{}
}

func (p *HelmDeploy) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	return mapstructure.Decode(pluginConfig, p)
}

func (p *HelmDeploy) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider) error {

	return nil
}

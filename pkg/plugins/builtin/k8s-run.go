package builtin

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/sirupsen/logrus"
)

type K8sRun struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	Name          string
	Targets       []K8sRunTargets
}

type K8sRunTargets struct {
	Name    string
	Secrets []interface{}
}

func (p *K8sRun) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	return mapstructure.Decode(pluginConfig, p)
}

func (p *K8sRun) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider) error {

	return nil
}

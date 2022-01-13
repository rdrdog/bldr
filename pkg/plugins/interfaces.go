package plugins

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/sirupsen/logrus"
)

type PluginDefinition interface {
	Execute(contextProvider *contexts.ContextProvider) error
	SetConfig(logger *logrus.Logger, targetName string, configuration *config.Configuration, pluginConfig map[string]interface{}) error
}

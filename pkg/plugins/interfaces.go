package plugins

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/sirupsen/logrus"
)

type PluginDefinition interface {
	Execute(contextProvider *contexts.ContextProvider, extensionsProvider *extensions.ExtensionsProvider) error
	SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error
}

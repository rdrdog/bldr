package plugins

import (
	"github.com/Redgwell/bldr/pkg/contexts"
	"github.com/sirupsen/logrus"
)

type PluginDefinition interface {
	Execute(contextProvider *contexts.ContextProvider) error
	SetConfig(logger *logrus.Logger, targetName string, pluginConfig map[string]interface{}) error
}

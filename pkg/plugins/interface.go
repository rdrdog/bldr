package plugins

import "github.com/sirupsen/logrus"

type PluginDefinition interface {
	Execute() error
	SetConfig(logger *logrus.Logger, pluginConfig map[string]interface{}) error
}

package plugins

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/sirupsen/logrus"
)

//counterfeiter:generate . PluginDefinition
type PluginDefinition interface {
	Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider, libProvider lib.LibProvider) error
	SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error
}

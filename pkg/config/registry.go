package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

// Build in packages are referred to without the fully qualified package path
const builtInPluginPathPrefix = "github.com/rdrdog/bldr/pkg/plugins/"
const builtInExtensionsPathPrefix = "github.com/rdrdog/bldr/pkg/extensions/"

type Registry struct {
	logger  *logrus.Logger
	plugins map[string]reflect.Type
}

func NewRegistry(logger *logrus.Logger) *Registry {
	registry := &Registry{
		logger:  logger,
		plugins: make(map[string]reflect.Type),
	}

	return registry
}

func (r *Registry) RegisterType(typeInterface interface{}) {
	t := reflect.TypeOf(typeInterface).Elem()

	// For internal packages and extensions, remove the fully qualified component:
	pluginName := strings.TrimPrefix(t.PkgPath(), builtInPluginPathPrefix)
	pluginName = strings.TrimPrefix(pluginName, builtInExtensionsPathPrefix)
	pluginName += "/" + t.Name()

	r.logger.Infof("🔌 adding plugin: %s", pluginName)
	r.plugins[pluginName] = t
}

func (r *Registry) CreateInstance(name string) (interface{}, error) {
	if r.plugins[name] == nil {
		return nil, fmt.Errorf("no plugin registered with name %s", name)
	}

	return reflect.New(r.plugins[name]).Interface(), nil
}

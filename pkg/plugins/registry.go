package plugins

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Redgwell/bldr/pkg/plugins/builtin"
	"github.com/sirupsen/logrus"
)

// Build in packages are referred to without the fully qualified package path
const packagePathPrefix = "github.com/Redgwell/bldr/pkg/plugins/"

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
	pluginName := strings.TrimPrefix(t.PkgPath(), packagePathPrefix) + "/" + t.Name()
	r.logger.Infof("Adding plugin: %s", pluginName)
	r.plugins[pluginName] = t
}

func (r *Registry) CreateInstance(name string) (PluginDefinition, error) {
	if r.plugins[name] == nil {
		return nil, fmt.Errorf("no plugin registered with name %s", name)
	}

	return reflect.New(r.plugins[name]).Interface().(PluginDefinition), nil
}

func (r *Registry) RegisterBuiltIn() {
	r.RegisterType((*builtin.BuildPathContextLoader)(nil))
	r.RegisterType((*builtin.DockerBuild)(nil))
	r.RegisterType((*builtin.GitContextLoader)(nil))
}

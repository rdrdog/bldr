package pipeline

import (
	"fmt"
	"reflect"
	"time"

	"github.com/rdrdog/bldr/internal/models"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	builtinExtensions "github.com/rdrdog/bldr/pkg/extensions/builtin"
	"github.com/rdrdog/bldr/pkg/plugins"
	builtinPlugins "github.com/rdrdog/bldr/pkg/plugins/builtin"
	"github.com/sirupsen/logrus"
)

type PluginPipeline struct {
	config             *config.Configuration
	contextProvider    *contexts.ContextProvider
	extensionsProvider *extensions.ExtensionsProvider
	logger             *logrus.Logger
	mode               string
	plugins            []plugins.PluginDefinition
	registry           *config.Registry
}

func (p *PluginPipeline) addPlugin(plugin plugins.PluginDefinition) {
	p.plugins = append(p.plugins, plugin)
}

// TODO - move this
func registerBuiltInPlugins(registry *config.Registry) {

	registry.RegisterType((*builtinPlugins.BuildPathContextLoader)(nil))
	registry.RegisterType((*builtinPlugins.DockerBuild)(nil))
	registry.RegisterType((*builtinPlugins.GitContextLoader)(nil))
	registry.RegisterType((*builtinPlugins.ManifestWriter)(nil))
	registry.RegisterType((*builtinExtensions.NullSecretLoader)(nil))

}

func NewPluginPipeline(logger *logrus.Logger, baseConfig *config.Configuration, pipelineOperationMode string) *PluginPipeline {
	registry := config.NewRegistry(logger)
	registerBuiltInPlugins(registry)

	pipeline := &PluginPipeline{
		config:             baseConfig,
		contextProvider:    contexts.NewContextProvider(logger),
		extensionsProvider: extensions.NewExtensionsProvider(logger, baseConfig, registry),
		logger:             logger,
		mode:               pipelineOperationMode,
		registry:           registry,
	}

	return pipeline
}

func (p *PluginPipeline) LoadPipelineStages() error {

	pipelineCfg, err := models.LoadPipelineConfig(p.config.Pipeline.Path)
	if err != nil {
		// log/blow up
		return err
	}

	var stages []models.Stage

	switch p.mode {
	case config.PipelineOperationModeBuild:
		p.extensionsProvider.LoadExtensions(pipelineCfg.Build.Extensions)
		stages = pipelineCfg.Build.Stages
	case config.PipelineOperationModeDeploy:
		p.extensionsProvider.LoadExtensions(pipelineCfg.Deploy.Extensions)
		stages = pipelineCfg.Deploy.Stages
	default:
		return fmt.Errorf("unexpected operation mode: '%s'", p.mode)
	}

	for i, s := range stages {

		p.logger.Infof("initialising target: %v using %v\n", s.Name, s.Plugin)
		// Load the PluginDefinition using the plugin registry for now
		// Later, we could potentially support go plugins

		pluginInterface, err := p.registry.CreateInstance(s.Plugin)
		if err != nil {
			return err
		}

		pluginInstance := pluginInterface.(plugins.PluginDefinition)

		yamlPath := fmt.Sprintf("$.%s.stages[%d].params", p.mode, i)
		pluginConfig := pipelineCfg.LoadPluginConfig(yamlPath)

		err = pluginInstance.SetConfig(p.logger, p.config, pluginConfig)
		if err != nil {
			return err
		}
		p.addPlugin(pluginInstance)
	}

	return nil
}

func (p *PluginPipeline) Run() error {

	pipelineStart := time.Now()

	for _, plugin := range p.plugins {
		pluginName := reflect.TypeOf(plugin).Elem().Name()
		p.logger.Infof("🚀 running plugin %s", pluginName)
		start := time.Now()

		err := plugin.Execute(p.contextProvider, p.extensionsProvider)

		p.logger.Infof("⏳ plugin %s took %v seconds", pluginName, time.Since(start).Seconds())

		if err != nil {

			// log/blow up
			p.logger.Errorf("failed running %v: %v", plugin, err)
			return err
		}
	}

	p.logger.Infof("✅ pipeline took %v seconds", time.Since(pipelineStart).Seconds())

	return nil
}

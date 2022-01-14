package pipeline

import (
	"fmt"
	"reflect"
	"time"

	"github.com/rdrdog/bldr/internal/models"
	"github.com/rdrdog/bldr/internal/providers"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/plugins"
	"github.com/sirupsen/logrus"
)

type PluginPipeline struct {
	config             *config.Configuration
	contextProvider    *providers.DefaultContextProvider
	extensionsProvider *providers.DefaultExtensionsProvider
	logger             *logrus.Logger
	mode               string
	plugins            []plugins.PluginDefinition
	registry           *providers.Registry
}

func (p *PluginPipeline) addPlugin(plugin plugins.PluginDefinition) {
	p.plugins = append(p.plugins, plugin)
}

func NewPluginPipeline(logger *logrus.Logger, baseConfig *config.Configuration, pipelineOperationMode string) *PluginPipeline {
	registry := providers.NewRegistry(logger)

	pipeline := &PluginPipeline{
		config:             baseConfig,
		contextProvider:    providers.NewContextProvider(logger),
		extensionsProvider: providers.NewExtensionsProvider(logger, baseConfig, registry),
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

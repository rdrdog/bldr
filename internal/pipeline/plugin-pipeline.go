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
	libProvider        *providers.DefaultLibProvider
	logger             *logrus.Logger
	plugins            []plugins.PluginDefinition
	registry           *providers.Registry
}

func (p *PluginPipeline) addPlugin(plugin plugins.PluginDefinition) {
	p.plugins = append(p.plugins, plugin)
}

func NewPluginPipeline(logger *logrus.Logger, cfg *config.Configuration) *PluginPipeline {
	registry := providers.NewRegistry(logger)

	pipeline := &PluginPipeline{
		config:             cfg,
		contextProvider:    providers.NewContextProvider(logger, cfg),
		extensionsProvider: providers.NewExtensionsProvider(logger, cfg, registry),
		libProvider:        providers.NewDefaultLibProvider(logger, cfg),
		logger:             logger,
		registry:           registry,
	}

	return pipeline
}

func (p *PluginPipeline) LoadPipelineStages(pipelineMode string) {

	pipelineCfg, err := models.LoadPipelineConfig(p.config.Paths.PipelineConfigFile)
	if err != nil {
		p.logger.Fatalf("could not load pipeline configuration: %v", err)
	}

	var stages []models.Stage

	switch pipelineMode {
	case config.PipelineOperationModeBuild:
		p.extensionsProvider.LoadExtensions(pipelineCfg.Build.Extensions)
		stages = pipelineCfg.Build.Stages
	case config.PipelineOperationModeDeploy:
		p.extensionsProvider.LoadExtensions(pipelineCfg.Deploy.Extensions)
		stages = pipelineCfg.Deploy.Stages
	default:
		p.logger.Fatalf("unexpected operation mode: '%s'", pipelineMode)
	}

	for i, s := range stages {

		p.logger.Infof("initialising target: %v using %v\n", s.Name, s.Plugin)
		// Load the PluginDefinition using the plugin registry for now
		// Later, we could potentially support go plugins

		pluginInterface, err := p.registry.CreateInstance(s.Plugin)
		if err != nil {
			p.logger.Fatalf("could not create plugin registry instance: %v", err)
		}

		pluginInstance := pluginInterface.(plugins.PluginDefinition)

		yamlPath := fmt.Sprintf("$.%s.stages[%d].params", pipelineMode, i)
		pluginConfig := pipelineCfg.LoadPluginConfig(yamlPath)

		err = pluginInstance.SetConfig(p.logger, p.config, pluginConfig)
		if err != nil {
			p.logger.Fatalf("could not initialise configuration for plugin %s on stage %s: %v", s.Plugin, s.Name, err)
		}
		p.addPlugin(pluginInstance)
	}
}

func (p *PluginPipeline) Run() error {

	pipelineStart := time.Now()

	for _, plugin := range p.plugins {
		pluginName := reflect.TypeOf(plugin).Elem().Name()
		p.logger.Infof("🚀 running plugin %s", pluginName)
		start := time.Now()

		err := plugin.Execute(p.contextProvider, p.extensionsProvider, p.libProvider)

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

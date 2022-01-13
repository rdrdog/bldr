package pipeline

import (
	"fmt"
	"reflect"
	"time"

	"github.com/rdrdog/bldr/internal/models"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/plugins"
	"github.com/sirupsen/logrus"
)

type PluginPipeline struct {
	config          *config.Configuration
	contextProvider *contexts.ContextProvider
	logger          *logrus.Logger
	mode            string
	plugins         []plugins.PluginDefinition
	registry        *plugins.Registry
}

func (p *PluginPipeline) addPlugin(plugin plugins.PluginDefinition) {
	p.plugins = append(p.plugins, plugin)
}

func NewPluginPipeline(logger *logrus.Logger, baseConfig *config.Configuration, mode string) *PluginPipeline {
	registry := plugins.NewRegistry(logger)
	registry.RegisterBuiltIn()

	pipeline := &PluginPipeline{
		config:          baseConfig,
		contextProvider: contexts.NewContextProvider(logger),
		logger:          logger,
		mode:            mode,
		registry:        registry,
	}
	return pipeline
}

func (p *PluginPipeline) AddPipelineConfigTargets() error {

	pipelineCfg, err := models.LoadPipelineConfig(p.config.Pipeline.Path)
	if err != nil {
		// log/blow up
		return err
	}

	for i, t := range pipelineCfg.Targets {
		p.logger.Infof("initialising target: %v using %v\n", t.Name, t.Build.Plugin)
		// Load the PluginDefinition using the plugin registry for now
		// Later, we could potentially support go plugins

		pluginInstance, err := p.registry.CreateInstance(t.Build.Plugin)
		if err != nil {
			return err
		}

		yamlPath := fmt.Sprintf("$.targets[%d].%s", i, p.mode)
		pluginConfig := pipelineCfg.LoadPluginConfig(yamlPath)

		err = pluginInstance.SetConfig(p.logger, t.Name, p.config, pluginConfig)
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

		err := plugin.Execute(p.contextProvider)

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

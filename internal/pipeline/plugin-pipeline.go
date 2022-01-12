package pipeline

import (
	"fmt"

	"github.com/Redgwell/bldr/internal/config"
	"github.com/Redgwell/bldr/pkg/plugins"
	"github.com/sirupsen/logrus"
)

type PluginPipeline struct {
	config  *config.Base
	logger  *logrus.Logger
	mode    string
	plugins []plugins.PluginDefinition
}

func NewPluginPipeline(logger *logrus.Logger, baseConfig *config.Base) *PluginPipeline {
	pipeline := &PluginPipeline{
		config: baseConfig,
		logger: logger,
		mode:   "build",
	}
	return pipeline
}

func (p *PluginPipeline) LoadPlugins(registry *plugins.Registry) error {

	pipelineCfg, err := config.LoadPipelineConfig(p.config.PipelineConfigPath)
	if err != nil {
		// log/blow up
		return err
	}

	for i, t := range pipelineCfg.Targets {
		p.logger.Printf("Plugin: %v - %v\n", t.Name, t.Build.Plugin)
		// Load the PluginDefinition using the plugin registry for now
		// Later, we could potentially support go plugins

		pluginInstance, err := registry.CreateInstance(t.Build.Plugin)
		if err != nil {
			return err
		}

		path := fmt.Sprintf("$.targets[%d].%s", i, p.mode)
		pluginConfig := pipelineCfg.LoadPluginConfig(path)
		p.logger.Printf("Plugin config for path %s: %v\n", path, pluginConfig)
		err = pluginInstance.SetConfig(p.logger, pluginConfig)
		if err != nil {
			return err
		}
		p.addPlugin(pluginInstance)
	}

	// // TODO - iterate through config, adding plugins?
	// p.addPlugin(&builtin.DockerBuild{})

	return nil
}

func (p *PluginPipeline) addPlugin(plugin plugins.PluginDefinition) {
	p.plugins = append(p.plugins, plugin)
}

func (p *PluginPipeline) Run() error {

	for _, plugin := range p.plugins {
		err := plugin.Execute()

		if err != nil {
			// log/blow up
			return err
		}
	}

	return nil
}

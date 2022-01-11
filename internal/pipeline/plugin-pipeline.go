package pipeline

import (
	"fmt"

	"github.com/Redgwell/bldr/internal/config"
	"github.com/Redgwell/bldr/pkg/plugins"
	"github.com/Redgwell/bldr/pkg/plugins/builtin"
)

/*
- Uses configuration to generate the plugin pipeline
- On execute, will run each phase of the pipeline
*/

type PluginPipeline struct {
	plugins []plugins.PluginDefinition
}

func NewPluginPipeline( /*todo - config, logger, etc*/ ) (*PluginPipeline, error) {
	pipeline := &PluginPipeline{}

	cfg, err := config.LoadPipelineConfig("samples/pipeline-config.yaml")
	if err != nil {
		// log/blow up
		return nil, err
	}

	for i, t := range cfg.Targets {
		fmt.Printf("Plugin: %v - %v\n", t.Name, t.Build.Plugin)

		path := fmt.Sprintf("$.targets[%d].deploy", i)
		pluginConfig := cfg.LoadPluginConfig(path)
		fmt.Printf("Plugin config for path %s: %v\n", path, pluginConfig)
		//fmt.Sprintln(pluginConfig)
	}

	// TODO - iterate through config, adding plugins?
	pipeline.addPlugin(&builtin.DockerBuild{})

	return pipeline, nil
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

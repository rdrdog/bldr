package pipeline

import "github.com/rdrdog/bldr/pkg/plugins/builtin"

func (p *PluginPipeline) AddDefaultPreBuildTargets() {
	buildPathContextLoader := &builtin.BuildPathContextLoader{}
	buildPathContextLoader.SetConfig(p.logger, "Build path context loader", p.config, nil)
	p.addPlugin(buildPathContextLoader)

	gitContextLoader := &builtin.GitContextLoader{}
	gitContextLoader.SetConfig(p.logger, "Git context loader", p.config, nil)
	p.addPlugin(gitContextLoader)
}

func (p *PluginPipeline) AddDefaultPostBuildTargets() {
	manifestWriter := &builtin.ManifestWriter{}
	manifestWriter.SetConfig(p.logger, "Manifest writer", p.config, nil)
	p.addPlugin(manifestWriter)
}

func (p *PluginPipeline) AddDefaultPreDeployTargets() {

}

func (p *PluginPipeline) AddDefaultPostDeployTargets() {

}

package pipeline

import "github.com/rdrdog/bldr/pkg/plugins/builtin"

func (p *PluginPipeline) AddDefaultPreBuildStages() {
	buildPathContextLoader := &builtin.BuildPathContextLoader{}
	buildPathContextLoader.SetConfig(p.logger, p.config, nil)
	p.addPlugin(buildPathContextLoader)

	gitContextLoader := &builtin.GitContextLoader{}
	gitContextLoader.SetConfig(p.logger, p.config, nil)
	p.addPlugin(gitContextLoader)
}

func (p *PluginPipeline) AddDefaultPostBuildStages() {
	manifestWriter := &builtin.ManifestWriter{}
	manifestWriter.SetConfig(p.logger, p.config, nil)
	p.addPlugin(manifestWriter)
}

func (p *PluginPipeline) AddDefaultPreDeployStages() {

}

func (p *PluginPipeline) AddDefaultPostDeployStages() {

}

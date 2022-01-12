package pipeline

import "github.com/Redgwell/bldr/pkg/plugins/builtin"

func (p *PluginPipeline) AddDefaultPreBuildTargets() error {
	buildPathContextLoader := &builtin.BuildPathContextLoader{}
	buildPathContextLoader.SetConfig(p.logger, "Build path context loader", p.config, nil)
	p.addPlugin(buildPathContextLoader)

	gitContextLoader := &builtin.GitContextLoader{}
	gitContextLoader.SetConfig(p.logger, "Git context loader", p.config, nil)
	p.addPlugin(gitContextLoader)

	return nil
}

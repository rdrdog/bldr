package cmd

import (
	"github.com/rdrdog/bldr/internal/pipeline"
	"github.com/rdrdog/bldr/pkg/config"
)

func PopulateBuildPipeline(pp *pipeline.PluginPipeline) {
	pp.AddDefaultPreBuildStages()

	pp.LoadPipelineStages(config.PipelineOperationModeBuild)

	pp.AddDefaultPostBuildStages()
}

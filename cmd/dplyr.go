package cmd

import (
	"github.com/rdrdog/bldr/internal/pipeline"
	"github.com/rdrdog/bldr/pkg/config"
)

func PopulateDeployPipeline(pp *pipeline.PluginPipeline, envName string) {
	pp.AddDefaultPreDeployStages(envName)

	pp.LoadPipelineStages(config.PipelineOperationModeDeploy)

	pp.AddDefaultPostDeployStages()
}

package main

import (
	"github.com/rdrdog/bldr/internal/pipeline"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	conf, err := config.Load(logger)
	if err != nil {
		logger.Fatalf("could no initialise configuration: %v", err)
		return
	}

	pp := pipeline.NewPluginPipeline(logger, conf, config.PipelineOperationModeBuild)

	pp.AddDefaultPreBuildTargets()

	err = pp.AddPipelineConfigTargets()
	if err != nil {
		logger.Fatalf("error loading plugins: %v", err)
	}

	pp.AddDefaultPostBuildTargets()

	pp.Run()
}

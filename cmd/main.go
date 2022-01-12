package main

import (
	"github.com/Redgwell/bldr/internal/pipeline"
	"github.com/Redgwell/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	config, err := config.Load(logger)
	if err != nil {
		logger.Fatalf("could no initialise configuration: %v", err)
		return
	}

	pp := pipeline.NewPluginPipeline(logger, config, "build")

	pp.AddDefaultPreBuildTargets()
	err = pp.AddPipelineConfigTargets()
	if err != nil {
		logger.Fatalf("error loading plugins: %v", err)
	}
	//pp.AddDefaultPostBuildTargets()
	pp.Run()
}

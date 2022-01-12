package main

import (
	"github.com/Redgwell/bldr/internal/config"
	"github.com/Redgwell/bldr/internal/pipeline"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	config, err := config.Load(logger)
	if err != nil {
		logger.Fatal("could no initialise configuration")
		return
	}

	pp := pipeline.NewPluginPipeline(logger, config, "build")

	pp.AddDefaultPreBuildTargets()
	err = pp.AddPipelineConfigTargets()
	if err != nil {
		logger.Fatalf("error loading plugins: %v", err)
	}
	pp.Run()
}

/*
TODO:
- main entrypoint to load config
	- Use env override stuff from WP
	- load config from .bldr
	- use default plugin suite for the command ('build', or 'deploy')
- Init the pluginPipeline from here then execute it
*/

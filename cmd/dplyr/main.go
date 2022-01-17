package main

import (
	"flag"

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

	// Parse args
	var envName string
	flag.StringVar(&envName, "e", "localk8s", "Specify environment name to deploy to. Default = localk8s")
	flag.Parse()

	pp := pipeline.NewPluginPipeline(logger, conf, config.PipelineOperationModeDeploy)

	pp.AddDefaultPreDeployStages(envName)

	err = pp.LoadPipelineStages()
	if err != nil {
		logger.Fatalf("error loading pipeline stages: %v", err)
	}

	pp.AddDefaultPostDeployStages()

	pp.Run()
}

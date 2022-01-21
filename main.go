package main

import (
	"flag"

	"github.com/rdrdog/bldr/cmd"
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

	flag.Parse()

	rawArgs := flag.Args()
	if len(rawArgs) == 0 {
		logger.Fatalf("either 'build' or 'deploy' command must be specified")
	}

	pp := pipeline.NewPluginPipeline(logger, conf)

	switch rawArgs[0] {
	case "build":
		cmd.PopulateBuildPipeline(pp)
	case "deploy":
		var envName string
		flag.StringVar(&envName, "e", "localk8s", "Specify environment name to deploy to. Default = localk8s")
		flag.Parse()

		cmd.PopulateDeployPipeline(pp, envName)
	}

	pp.Run()
}

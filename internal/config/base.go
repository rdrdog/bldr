package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Base struct {
	CI       bool `env:"CI" envDefault:"false"`
	Docker   DockerConfig
	Logging  LoggingConfig
	Pipeline PipelineConfig
}

func Load(logger *logrus.Logger) (*Base, error) {
	newConfig := &Base{}

	for _, configSection := range []interface{}{
		newConfig,
		&newConfig.Docker,
		&newConfig.Logging,
		&newConfig.Pipeline,
	} {
		if err := env.Parse(configSection); err != nil {
			return nil, fmt.Errorf("unable to load the config: %v", err)
		}
	}

	// detect if we're running in CI, or local
	if newConfig.CI {
		logger.Info("Configuring for CI environment")
		newConfig.Docker.IncludeTimeInImageTag = false
		newConfig.Docker.UseRemoteContainerRegistryCache = false
	} else {
		logger.Info("Configuring for local environment")
		newConfig.Docker.IncludeTimeInImageTag = true
		newConfig.Docker.UseRemoteContainerRegistryCache = true
		newConfig.Docker.Registry = ""
		newConfig.Pipeline.Path = "samples/pipeline-config.yaml"
		newConfig.Logging.Level = "DEBUG"
	}

	newConfig.Logging.SetFormatter(logger)

	return newConfig, nil
}

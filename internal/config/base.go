package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Base struct {
	Logging            LoggingConfig
	PipelineConfigPath string
}

func Load(logger *logrus.Logger) (*Base, error) {
	newConfig := &Base{
		PipelineConfigPath: "samples/pipeline-config.yaml",
	}

	for _, configSection := range []interface{}{
		newConfig,
		&newConfig.Logging,
	} {
		if err := env.Parse(configSection); err != nil {
			return nil, fmt.Errorf("unable to load the config: %v", err)
		}
	}

	newConfig.Logging.SetFormatter(logger)

	return newConfig, nil
}

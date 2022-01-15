package builtin

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

type LocalSecretLoader struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	Path          string
}

func (e *LocalSecretLoader) SetConfig(logger *logrus.Logger, configuration *config.Configuration, extensionConfig map[string]interface{}) error {
	e.configuration = configuration
	e.logger = logger
	err := mapstructure.Decode(extensionConfig, e)
	if err != nil {
		return err
	}

	if len(e.Path) == 0 {
		e.logger.Infof("Path not specified - defaulting local secrets to ")
	}
	return nil
}

func (e *LocalSecretLoader) LoadSecrets(map[string]interface{}) map[string]string {
	return nil
}

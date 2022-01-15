package builtin

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

type NullSecretLoader struct {
}

func (e *NullSecretLoader) SetConfig(logger *logrus.Logger, configuration *config.Configuration, extensionConfig map[string]interface{}) error {
	return nil
}

func (e *NullSecretLoader) LoadSecrets(map[string]interface{}) map[string]string {
	return nil
}

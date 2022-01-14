package extensions

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

type ExtensionsProvider interface {
	GetSecretLoader() SecretLoader
}
type SecretLoader interface {
	SetConfig(logger *logrus.Logger, configuration *config.Configuration, extensionConfig map[string]interface{}) error
	LoadSecrets(map[string]interface{}) map[string]string
}

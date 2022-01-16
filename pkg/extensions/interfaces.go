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
	LoadSecrets(targetName string, secretParams []interface{}) ([]*SecretKeyValuePair, error)
}

type SecretKeyValuePair struct {
	Key   string
	Value string
}

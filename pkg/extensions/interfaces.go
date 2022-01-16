package extensions

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

//counterfeiter:generate . ExtensionsProvider
type ExtensionsProvider interface {
	GetSecretLoader() SecretLoader
}

//counterfeiter:generate . SecretLoader
type SecretLoader interface {
	SetConfig(logger *logrus.Logger, configuration *config.Configuration, extensionConfig map[string]interface{}) error
	LoadSecrets(targetName string, secretParams []interface{}) ([]*SecretKeyValuePair, error)
}

type SecretKeyValuePair struct {
	Key   string
	Value string
}

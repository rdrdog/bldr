package extensions

import (
	"strings"

	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/extensions/builtin"
	"github.com/sirupsen/logrus"
)

type ExtensionsProvider struct {
	SecretLoader  SecretsLoader
	configuration *config.Configuration
	logger        *logrus.Logger
	registry      *config.Registry
}

type ExtensionDefinition struct {
	Definition string
	Params     map[string]interface{}
}

func NewExtensionsProvider(logger *logrus.Logger, configuration *config.Configuration, registry *config.Registry) *ExtensionsProvider {

	return &ExtensionsProvider{
		configuration: configuration,
		logger:        logger,
		registry:      registry,
		SecretLoader:  &builtin.NullSecretLoader{},
	}
}

func (e *ExtensionsProvider) LoadExtensions(extensions map[string]ExtensionDefinition) error {
	for key, value := range extensions {

		e.logger.Infof("Loading extension %s as %s", key, value.Definition)

		switch strings.ToLower(key) {
		case "secretloader":
			instance, err := e.registry.CreateInstance(value.Definition)
			if err != nil {
				e.logger.Errorf("failed to create instance of %s: %v", value.Definition, err)
				return err
			}

			e.SecretLoader = instance.(SecretsLoader)
			err = e.SecretLoader.SetConfig(e.logger, e.configuration, value.Params)
			if err != nil {
				e.logger.Errorf("failed to load secret loader extension: %v", err)
				return err
			}
		}
	}

	return nil
}

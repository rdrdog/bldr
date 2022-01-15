package providers

import (
	"strings"

	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/extensions/builtin"
	"github.com/sirupsen/logrus"
)

type DefaultExtensionsProvider struct {
	secretLoader  extensions.SecretLoader
	configuration *config.Configuration
	logger        *logrus.Logger
	registry      *Registry
}

type ExtensionDefinition struct {
	Definition string
	Params     map[string]interface{}
}

func NewExtensionsProvider(logger *logrus.Logger, configuration *config.Configuration, registry *Registry) *DefaultExtensionsProvider {

	return &DefaultExtensionsProvider{
		configuration: configuration,
		logger:        logger,
		registry:      registry,
		secretLoader:  &builtin.LocalSecretLoader{},
	}
}

func (e *DefaultExtensionsProvider) LoadExtensions(defs map[string]ExtensionDefinition) error {
	for key, value := range defs {

		e.logger.Infof("Loading extension %s as %s", key, value.Definition)

		switch strings.ToLower(key) {
		case "secretloader":
			instance, err := e.registry.CreateInstance(value.Definition)
			if err != nil {
				e.logger.Errorf("failed to create instance of %s: %v", value.Definition, err)
				return err
			}

			e.secretLoader = instance.(extensions.SecretLoader)
			err = e.secretLoader.SetConfig(e.logger, e.configuration, value.Params)
			if err != nil {
				e.logger.Errorf("failed to load secret loader extension: %v", err)
				return err
			}
		}
	}

	return nil
}

func (p *DefaultExtensionsProvider) GetSecretLoader() extensions.SecretLoader {
	return p.secretLoader
}

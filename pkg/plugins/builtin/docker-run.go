package builtin

import (
	"github.com/gookit/goutil/arrutil"
	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/sirupsen/logrus"
)

type DockerRun struct {
	logger           *logrus.Logger
	Name             string
	SkipEnvironments []string
	Targets          []DockerRunTargets
}

type DockerRunTargets struct {
	Name    string
	Secrets []interface{}
}

func (p *DockerRun) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.logger = logger
	return mapstructure.Decode(pluginConfig, p)
}

func (p *DockerRun) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider, libProvider lib.LibProvider) error {
	dc := contextProvider.GetDeployContext()

	if arrutil.Contains(p.SkipEnvironments, dc.EnvironmentName) {
		p.logger.Infof("⏩ %s is listed in skip environments - skipping docker run", dc.EnvironmentName)
		return nil
	}

	p.logger.Debugf("environment '%s' is not listed in skip environments %v - running targets", dc.EnvironmentName, p.SkipEnvironments)

	docker := libProvider.GetDockerLib()

	for _, t := range p.Targets {
		secrets, err := extensionsProvider.GetSecretLoader().LoadSecrets(t.Name, t.Secrets)
		if err != nil {
			p.logger.Errorf("error loading secrets for target %s", t.Name)
			return err
		}
		p.logger.Debugf("loaded %d secrets for target %s", len(secrets), t.Name)

		imageNameAndTag := dc.GetArtefactByName(t.Name)
		err = docker.RunImage(imageNameAndTag, secretsToMap(secrets), nil)

		if err != nil {
			return err
		}
	}

	return nil
}

func secretsToMap(s []*extensions.SecretKeyValuePair) map[string]string {
	result := make(map[string]string, len(s))
	for _, val := range s {
		result[val.Key] = val.Value
	}

	return result
}

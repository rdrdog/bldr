package builtin

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type DockerBuild struct {
	logger  *logrus.Logger
	Path    string
	Include []string
}

func (p *DockerBuild) SetConfig(logger *logrus.Logger, pluginConfig map[string]interface{}) error {
	p.logger = logger
	return mapstructure.Decode(pluginConfig, p)
}

func (p *DockerBuild) Execute( /*projectName string, targetName string*/ ) error {
	p.logger.Infof("Running docker build with config: Path: %s, Include: %d", p.Path, p.Include)
	return nil
}

package builtin

import (
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const ManifestFileName = "manifest.yaml"

type ManifestWriter struct {
	config *config.Configuration
	logger *logrus.Logger
	Name   string
}

func (p *ManifestWriter) SetConfig(logger *logrus.Logger, targetName string, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.config = configuration
	p.logger = logger
	p.Name = targetName
	return mapstructure.Decode(pluginConfig, p)
}

func (p *ManifestWriter) Execute(contextProvider *contexts.ContextProvider) error {

	bc := contextProvider.BuildContext
	manifestInstance := &manifest{
		BuildNumber: bc.BuildNumber,
		Artefacts:   bc.ArtefactManifest.Artefacts,
		MetaData: metaData{
			BldrVersion:     config.BldrAppVersion,
			ManifestVersion: config.ManifestVersion,
		},
	}
	manifestInstance.Repo.BranchName = bc.GitContext.BranchName
	manifestInstance.Repo.CommitSha = bc.GitContext.FullCommitSha

	manifestData, err := yaml.Marshal(manifestInstance)
	p.logger.Debugf("Writing manifest data: \n%v", string(manifestData))

	if err != nil {
		p.logger.Errorf("error generating manifest yaml: %v", err)
		return err
	}

	manifestFilePath := path.Join(bc.PathContext.BuildArtefactDirectory, ManifestFileName)
	err = afero.WriteFile(config.Appfs, manifestFilePath, manifestData, 0777)
	if err != nil {
		p.logger.Errorf("error writing manifest file to %s: %v", manifestFilePath, err)
		return err
	}

	return nil
}

type manifest struct {
	BuildNumber string
	Repo        struct {
		BranchName string
		CommitSha  string
	}
	Artefacts map[string]string
	MetaData  metaData
}

type metaData struct {
	BldrVersion     string
	ManifestVersion string
}

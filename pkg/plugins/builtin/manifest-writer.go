package builtin

import (
	"path"

	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/rdrdog/bldr/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const ManifestFileName = "manifest.yaml"

type ManifestWriter struct {
	config *config.Configuration
	logger *logrus.Logger
}

func (p *ManifestWriter) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.config = configuration
	p.logger = logger
	return nil
}

func (p *ManifestWriter) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider, libProvider lib.LibProvider) error {

	bc := contextProvider.GetBuildContext()
	manifestInstance := &models.Manifest{
		BuildNumber: bc.BuildNumber,
		Artefacts:   bc.ArtefactManifest.Artefacts,
		MetaData: models.ManifestMetaData{
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

	manifestFilePath := path.Join(p.config.Paths.BuildArtefactDirectory, ManifestFileName)
	err = afero.WriteFile(config.Appfs, manifestFilePath, manifestData, 0755)
	if err != nil {
		p.logger.Errorf("error writing manifest file to %s: %v", manifestFilePath, err)
		return err
	}

	return nil
}

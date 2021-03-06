package builtin

import (
	"testing"

	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/models"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestManifestWriter_SetConfig(t *testing.T) {
	mw := &ManifestWriter{}

	err := mw.SetConfig(nil, nil, nil)

	assert.Nil(t, err)

}

func TestManifestWriter_Execute_File_Exists_And_Has_Data(t *testing.T) {

	mw := &ManifestWriter{}
	logger, _ := test.NewNullLogger()
	mw.logger = logger
	mw.config = &config.Configuration{}
	mw.config.Paths = config.PathsConfig{
		RepoRootDirectory:      ".",
		BuildArtefactDirectory: ".",
	}

	config.Appfs = afero.NewMemMapFs()

	buildContext := &contexts.BuildContext{
		BuildNumber: "99999",
		ArtefactManifest: &contexts.ArtefactManifest{
			Artefacts: make(map[string]string),
		},
		GitContext: &contexts.GitContext{
			FullCommitSha:  "8f825b7454e61e5c61fcf165e8299151b46d67f8",
			ShortCommitSha: "8f825b7",
			BranchName:     "main",
		},
	}
	contextProvider := &mockContextProvider{
		buildContext: buildContext,
	}

	mw.Execute(contextProvider, nil, nil)

	fileExists, _ := afero.Exists(config.Appfs, "./manifest.yaml")
	// Check that the file exists
	assert.True(t, fileExists)

	// Check that the data in the file is the same as the one passed in by the test.
	data, _ := afero.ReadFile(config.Appfs, "manifest.yaml")

	var result models.Manifest
	yaml.Unmarshal(data, &result)

	assert.Equal(t, buildContext.BuildNumber, result.BuildNumber)
	assert.Equal(t, buildContext.GitContext.FullCommitSha, result.Repo.CommitSha)
	assert.Equal(t, buildContext.GitContext.BranchName, result.Repo.BranchName)
	assert.Equal(t, buildContext.ArtefactManifest.Artefacts, result.Artefacts)
}

type mockContextProvider struct {
	buildContext *contexts.BuildContext
}

func (p *mockContextProvider) GetBuildContext() *contexts.BuildContext {
	return p.buildContext
}

func (p *mockContextProvider) GetDeployContext() *contexts.DeployContext {
	return nil
}

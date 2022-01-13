package builtin

import (
	"testing"

	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestManifestWriter_SetConfig(t *testing.T) {
	mw := &ManifestWriter{}

	err := mw.SetConfig(nil, "test", nil, nil)

	assert.Nil(t, err)
	assert.Equal(t, "test", mw.Name)

}

func TestManifestWriter_Execute_File_Exists_And_Has_Data(t *testing.T) {

	mw := &ManifestWriter{}
	logger, _ := test.NewNullLogger()
	mw.logger = logger

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
		PathContext: &contexts.PathContext{
			RepoRootDirectory:      ".",
			BuildArtefactDirectory: ".",
		},
	}
	mw.Execute(&contexts.ContextProvider{
		BuildContext: buildContext,
	})

	fileExists, _ := afero.Exists(config.Appfs, "./manifest.yaml")
	// Check that the file exists
	assert.True(t, fileExists)

	// Check that the data in the file is the same as the one passed in by the test.
	data, _ := afero.ReadFile(config.Appfs, "manifest.yaml")

	var result manifest
	yaml.Unmarshal(data, &result)

	assert.Equal(t, buildContext.BuildNumber, result.BuildNumber)
	assert.Equal(t, buildContext.GitContext.FullCommitSha, result.Repo.CommitSha)
	assert.Equal(t, buildContext.GitContext.BranchName, result.Repo.BranchName)
	assert.Equal(t, buildContext.ArtefactManifest.Artefacts, result.Artefacts)
}

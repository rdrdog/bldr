package builtin

import (
	"path"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/contexts/contextsfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const testManifestFilePath = "./output/manifest.yaml"
const testManifestFileContents = `
artefacts:
  infra-example: infra-example:2bf0f93-220116091041
metadata:
  bldrversion: "1.0"
  manifestversion: "1.0"
`

func getConfiguredDeployContextLoader(t *testing.T) *DeployContextLoader {
	paramsYaml := `
environmentName: dev
`

	params := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(paramsYaml), params)
	assert.Nil(t, err)

	logger, _ := test.NewNullLogger()
	config := new(config.Configuration)
	config.Paths.DeploymentManifestFile = testManifestFilePath

	p := &DeployContextLoader{}
	p.SetConfig(logger, config, params)

	return p
}

func TestDeployContextLoaderSetConfig(t *testing.T) {
	p := getConfiguredDeployContextLoader(t)

	assert.Equal(t, "dev", p.EnvironmentName)
}

func TestDeployContextLoaderExecute(t *testing.T) {
	spec.Run(t, "DeployContextLoader.Execute", testDeployContextLoaderExecute, spec.Report(report.Terminal{}))
}

func testDeployContextLoaderExecute(t *testing.T, when spec.G, it spec.S) {
	var p *DeployContextLoader
	var fakeContextProvider *contextsfakes.FakeContextProvider
	var deployContext *contexts.DeployContext

	it.Before(func() {
		RegisterTestingT(t)
		config.Appfs = afero.NewMemMapFs()
		config.Appfs.MkdirAll(path.Dir(testManifestFilePath), 0755)
		afero.WriteFile(config.Appfs, testManifestFilePath, []byte(testManifestFileContents), 0755)

		p = getConfiguredDeployContextLoader(t)
		fakeContextProvider = &contextsfakes.FakeContextProvider{}

		deployContext = &contexts.DeployContext{
			Artefacts: make(map[string]string),
		}
		fakeContextProvider.GetDeployContextReturns(deployContext)
	})

	it("sets the deployment context environment name", func() {
		err := p.Execute(fakeContextProvider, nil, nil)
		assert.Nil(t, err)
		Expect(deployContext.EnvironmentName).To(Equal("dev"))
	})

	it("returns an error if the manifest file does not exist", func() {
		p.configuration.Paths.DeploymentManifestFile = "nothing.yaml"
		err := p.Execute(fakeContextProvider, nil, nil)
		assert.NotNil(t, err)
	})

	it("returns an error if the manifest file is not yaml", func() {
		afero.WriteFile(config.Appfs, testManifestFilePath, []byte("not yaml"), 0755)
		err := p.Execute(fakeContextProvider, nil, nil)
		assert.NotNil(t, err)
	})
}

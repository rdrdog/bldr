package builtin

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/contexts/contextsfakes"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/extensions/extensionsfakes"
	"github.com/rdrdog/bldr/pkg/lib/libfakes"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func getConfiguredDockerRun(t *testing.T) *DockerRun {
	paramsYaml := `
skipenvironments:
  - localk8s
targets:
  - name: infra-example
    secrets:
      - key: SOME_KEY
        envValue: SOME_ENV_VALUE
  - name: another-infra-example
`

	params := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(paramsYaml), params)
	assert.Nil(t, err)

	logger, _ := test.NewNullLogger()

	p := &DockerRun{}
	p.SetConfig(logger, nil, params)

	return p
}

func TestDockerRunSetConfig(t *testing.T) {
	p := getConfiguredDockerRun(t)

	assert.Contains(t, p.SkipEnvironments, "localk8s")
	assert.Equal(t, 2, len(p.Targets))
	assert.Equal(t, "infra-example", p.Targets[0].Name)
	assert.NotNil(t, p.Targets[0].Secrets)
}

func TestDockerRunExecute(t *testing.T) {
	spec.Run(t, "DockerRun.Execute", testDockerRunExecute, spec.Report(report.Terminal{}))
}

func testDockerRunExecute(t *testing.T, when spec.G, it spec.S) {
	var p *DockerRun
	var fakeContextProvider *contextsfakes.FakeContextProvider
	var fakeExtensionsProvider *extensionsfakes.FakeExtensionsProvider
	var fakeLibProvider *libfakes.FakeLibProvider
	var fakeSecretLoader *extensionsfakes.FakeSecretLoader
	var fakeDockerLib *libfakes.FakeDocker

	it.Before(func() {
		RegisterTestingT(t)
		p = getConfiguredDockerRun(t)
		fakeContextProvider = &contextsfakes.FakeContextProvider{}
		fakeExtensionsProvider = &extensionsfakes.FakeExtensionsProvider{}
		fakeLibProvider = &libfakes.FakeLibProvider{}
		fakeSecretLoader = &extensionsfakes.FakeSecretLoader{}
		fakeDockerLib = &libfakes.FakeDocker{}

		fakeExtensionsProvider.GetSecretLoaderReturns(fakeSecretLoader)
		fakeLibProvider.GetDockerLibReturns(fakeDockerLib)

		dc := &contexts.DeployContext{
			EnvironmentName: "dev",
			Artefacts:       make(map[string]string),
		}
		dc.Artefacts["infra-example"] = "infra-example:abc123"
		dc.Artefacts["another-infra-example"] = "another-infra-example:def456"
		fakeContextProvider.GetDeployContextReturns(dc)

	})

	it("skips deployment if the environment name is configured in the skipEnvironments", func() {
		fakeContextProvider.GetDeployContext().EnvironmentName = "localk8s"

		err := p.Execute(fakeContextProvider, fakeExtensionsProvider, fakeLibProvider)
		assert.Nil(t, err)
		Expect(fakeLibProvider.GetDockerLibCallCount()).To(Equal(0))
	})

	it("returns an error if secret loading fails", func() {
		fakeSecretLoader.LoadSecretsReturns(nil, errors.New("oops"))
		err := p.Execute(fakeContextProvider, fakeExtensionsProvider, fakeLibProvider)
		assert.NotNil(t, err)
	})

	it("gets secrets from the secret loader", func() {
		err := p.Execute(fakeContextProvider, fakeExtensionsProvider, fakeLibProvider)
		assert.Nil(t, err)
		Expect(fakeSecretLoader.LoadSecretsCallCount()).To(Equal(2))
	})

	it("calls docker.RunImage with the expected arguments", func() {
		secrets := []*extensions.SecretKeyValuePair{
			new(extensions.SecretKeyValuePair),
		}
		secrets[0].Key = "abc123"
		secrets[0].Value = "def456"

		secretsMap := map[string]string{
			"abc123": "def456",
		}

		fakeSecretLoader.LoadSecretsReturns(secrets, nil)

		err := p.Execute(fakeContextProvider, fakeExtensionsProvider, fakeLibProvider)
		assert.Nil(t, err)
		Expect(fakeDockerLib.RunImageCallCount()).To(Equal(2))
		runImageImageName, runImageSecrets, runImageBuildArgs := fakeDockerLib.RunImageArgsForCall(0)
		Expect(runImageImageName).To(Equal("infra-example:abc123"))
		Expect(runImageSecrets).To(Equal(secretsMap))
		Expect(runImageBuildArgs).To(BeNil())
	})

	it("returns an error if docker.RunImage fails", func() {
		secrets := []*extensions.SecretKeyValuePair{
			new(extensions.SecretKeyValuePair),
		}
		fakeSecretLoader.LoadSecretsReturns(secrets, nil)
		fakeDockerLib.RunImageReturns(errors.New("fail"))

		err := p.Execute(fakeContextProvider, fakeExtensionsProvider, fakeLibProvider)
		assert.NotNil(t, err)
		// Ensure the runImage stopped after the first failure
		Expect(fakeDockerLib.RunImageCallCount()).To(Equal(1))
	})
}

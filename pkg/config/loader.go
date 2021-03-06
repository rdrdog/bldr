package config

import (
	"fmt"
	"os"
	"path"

	"github.com/caarlos0/env"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const buildArtefactDirectoryName = "build-artefacts"
const buildEnvironmentNameLocal = "local"
const buildEnvironmentNameCI = "ci"
const bldrConfigFileName = "bldr.yaml"
const bldrConfigDefaults = `
default:
  docker:
    useBuildKit: true
  git:
    mainBranchName: main
  logging:
    level: INFO
  paths:
    deploymentManifestFile: build-artefacts/manifest.yaml
    pipelineConfigFile: samples/pipeline-config.yaml

local:
  docker:
    includeTimeInImageTag: true
    pushContainers: false
    registry: ""
    useRemoteContainerRegistryCache: false

ci:
  docker:
    includeTimeInImageTag: false
    pushContainers: true
    registry: "todo:set-your-container-registry"
    useRemoteContainerRegistryCache: true
`

func getConfigDefaults() *Configuration {
	configDefaults := &Configuration{
		Docker: DockerConfig{
			UseBuildKit: true,
		},
		Git: GitConfig{
			MainBranchName: "main",
		},
		Logging: LoggingConfig{
			Level: "INFO",
		},
		Paths: PathsConfig{
			DeploymentManifestFile: "build-artefacts/manifest.yaml",
			PipelineConfigFile:     "pipeline-config.yaml",
		},
	}

	return configDefaults
}

func Load(logger *logrus.Logger) (*Configuration, error) {
	newConfig := getConfigDefaults()

	// First, populate any environment var overrides
	for _, configSection := range []interface{}{
		newConfig,
		&newConfig.Docker,
		&newConfig.Git,
		&newConfig.Logging,
		&newConfig.Paths,
	} {
		if err := env.Parse(configSection); err != nil {
			return nil, fmt.Errorf("unable to load the config: %v", err)
		}
	}

	// Next, load our bldr.yaml settings
	bldrConfigData, err := loadOrGenerateBldrConfig(logger)
	if err != nil {
		return nil, fmt.Errorf("unable to load the %s config: %v", bldrConfigFileName, err)
	}

	// Populate our config using the bldr config
	mapstructure.Decode(bldrConfigData["default"], newConfig)
	if newConfig.CI {
		logger.Info("Configuring for CI environment")
		mapstructure.Decode(bldrConfigData[buildEnvironmentNameCI], newConfig)
	} else {
		logger.Info("Configuring for local environment")
		mapstructure.Decode(bldrConfigData[buildEnvironmentNameLocal], newConfig)
	}

	newConfig.Logging.Configure(logger)

	newConfig.Paths.RepoRootDirectory, err = os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to load current working directory: %v", err)
	}

	logger.Debugf("repo root directory determined to be %s", newConfig.Paths.RepoRootDirectory)
	// Set the build artefact directory to be rooted at the repo root directory
	newConfig.Paths.BuildArtefactDirectory = path.Join(newConfig.Paths.RepoRootDirectory, buildArtefactDirectoryName)

	return newConfig, nil
}

func loadOrGenerateBldrConfig(logger *logrus.Logger) (map[string]interface{}, error) {
	exists, err := afero.Exists(Appfs, bldrConfigFileName)
	if err != nil {
		logger.Errorf("error checking if bldr config file exists at %s: %v", bldrConfigFileName, err)
		return nil, err
	}

	// If our bldr config file doesn't exist, generate it
	if !exists {
		logger.Infof("no bldr config file found at %s - generating defaults", bldrConfigFileName)
		err = afero.WriteFile(Appfs, bldrConfigFileName, []byte(bldrConfigDefaults), 0777)
		if err != nil {
			logger.Errorf("error writing bldr config file to %s: %v", bldrConfigFileName, err)
			return nil, err
		}
	}

	// Read the bldr config:
	data, err := afero.ReadFile(Appfs, bldrConfigFileName)
	if err != nil {
		logger.Errorf("error reading bldr config file from %s: %v", bldrConfigFileName, err)
		return nil, err
	}

	// Map the yaml to a dictionary structure
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		logger.Errorf("error loading bldr config file yaml from %s: %v", bldrConfigFileName, err)
		return nil, err
	}

	return result, err
}

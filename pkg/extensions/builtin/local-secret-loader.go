package builtin

import (
	"path"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type LocalSecretLoader struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	Path          string
}

func (e *LocalSecretLoader) SetConfig(logger *logrus.Logger, configuration *config.Configuration, extensionConfig map[string]interface{}) error {
	e.configuration = configuration
	e.logger = logger
	err := mapstructure.Decode(extensionConfig, e)
	if err != nil {
		e.logger.Errorf("unable to map config onto localsecretloader: %v", err)
		return err
	}

	if len(e.Path) == 0 {
		e.Path = "secrets"
	}
	// If the configured path is not absolute, then set it to be inside the repo
	if !path.IsAbs(e.Path) {
		e.Path = path.Join(e.configuration.Paths.RepoRootDirectory, e.Path)
	}

	e.logger.Infof("local secrets path set to %s", e.Path)
	err = config.Appfs.MkdirAll(e.Path, 0775)
	if err != nil {
		e.logger.Errorf("unable to create secrets folder at %s: %v", e.Path, err)
		return err
	}

	return nil
}

/*
	- For each item in secretItem
	- mapstructure.Decode it to an array of secretConfig
	- check if the folder exists for the secret
		- the folder will be based on the targetName
		- e.g path.Join(e.Path, targetName)
	- check if a file exists for the secretItem.key
		- If not, create an empty file
	- return a map of secretItem.envValue = secretValue for each one
*/
func (e *LocalSecretLoader) LoadSecrets(targetName string, secretParams []interface{}) ([]*extensions.SecretKeyValuePair, error) {
	result := make([]*extensions.SecretKeyValuePair, len(secretParams))
	e.logger.Infof("fetching secrets for target %s", targetName)

	for i, val := range secretParams {
		secret := &secretItem{}
		err := mapstructure.Decode(val, secret)
		if err != nil {
			e.logger.Errorf("could not decode secret params at index %d for target %s: %v", i, targetName, err)
			return nil, err
		}

		secretFolder := path.Join(e.Path, targetName)
		dirExists, err := afero.DirExists(config.Appfs, secretFolder)
		if err != nil {
			e.logger.Errorf("could not determine if secret folder exists at %s: %v", secretFolder, err)
			return nil, err
		}

		if !dirExists {
			err = config.Appfs.MkdirAll(secretFolder, 0755)
			if err != nil {
				e.logger.Errorf("could not create secret folder at %s: %v", secretFolder, err)
				return nil, err
			}
		}

		secretFile := path.Join(secretFolder, secret.Key)
		e.logger.Debugf("loading secret from %s", secretFile)
		secretFileExists, err := afero.Exists(config.Appfs, secretFile)
		if err != nil {
			e.logger.Errorf("could not determine if secret file exists at %s: %v", secretFile, err)
			return nil, err
		}

		if !secretFileExists {
			err = afero.WriteFile(config.Appfs, secretFile, nil, 0755)
			if err != nil {
				e.logger.Errorf("could not create secret file at %s: %v", secretFile, err)
				return nil, err
			}
		}

		data, err := afero.ReadFile(config.Appfs, secretFile)
		if err != nil {
			e.logger.Errorf("could not read secret file contents from %s: %v", secretFile, err)
			return nil, err
		}

		value := strings.TrimSpace(string(data))
		// add the value to the logger's obfuscated values so it's not logged out
		e.configuration.Logging.AddToSecretMask(value)

		result[i] = &extensions.SecretKeyValuePair{
			Key:   secret.EnvValue,
			Value: value,
		}
	}

	return result, nil
}

type secretItem struct {
	Key      string // The key of the secret
	EnvValue string // The environment variable name to define the secret value as
}

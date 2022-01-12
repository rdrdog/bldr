package models

import (
	"bytes"
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

type PipelineConfig struct {
	Targets []Target
	source  []byte
}

type Target struct {
	Name   string
	Build  BuildTarget
	Deploy DeployTarget
}

type BuildTarget struct {
	Plugin  string
	Path    string
	Include []string
}

type DeployTarget struct {
	Plugin string
}

func LoadPipelineConfig(configFilePath string) (*PipelineConfig, error) {
	cfg := PipelineConfig{}

	var err error
	cfg.source, err = ioutil.ReadFile(configFilePath)
	if err != nil {
		// TODO log
		return nil, err
	}
	yaml.Unmarshal(cfg.source, &cfg)

	return &cfg, nil
}

func (pc *PipelineConfig) LoadPluginConfig(path string) map[string]interface{} {
	subset, _ := yaml.PathString(path)

	var result map[string]interface{}
	subset.Read(bytes.NewReader(pc.source), &result)

	return result
}

package config

type Configuration struct {
	CI      bool `env:"CI" envDefault:"false"`
	Docker  DockerConfig
	Git     GitConfig
	Logging LoggingConfig
	Paths   PathsConfig
}

type DockerConfig struct {
	Registry                        string
	IncludeTimeInImageTag           bool
	PushContainers                  bool
	UseBuildKit                     bool
	UseRemoteContainerRegistryCache bool
}

type GitConfig struct {
	MainBranchName string
}

type LoggingConfig struct {
	Level string
}

type PathsConfig struct {
	BuildArtefactDirectory string // not configurable, set in config/loader.go
	PipelineConfigFile     string
	RepoRootDirectory      string // not configurable, set in config/loader.go
}

package config

type Configuration struct {
	CI       bool `env:"CI" envDefault:"false"`
	Docker   DockerConfig
	Git      GitConfig
	Logging  LoggingConfig
	Pipeline PipelineConfig
}

type LoggingConfig struct {
	Level string
}

type PipelineConfig struct {
	Path string
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

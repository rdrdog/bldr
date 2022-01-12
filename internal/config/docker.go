package config

type DockerConfig struct {
	Registry              string `env:"DOCKER_REGISTRY" envDefault:""`
	IncludeTimeInImageTag bool   // default to true for local builds, false for cloud builds
}

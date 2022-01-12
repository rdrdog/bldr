package config

type PipelineConfig struct {
	Path string `env:"PIPELINE_CONFIG_PATH" envDefault:"pipeline-config.yaml"`
}

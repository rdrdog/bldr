package contexts

type DeployContext struct {
	EnvironmentName string
	Artefacts       map[string]string
}

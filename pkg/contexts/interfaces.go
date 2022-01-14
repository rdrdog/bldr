package contexts

type ContextProvider interface {
	GetBuildContext() *BuildContext
	GetDeployContext() *DeployContext
}

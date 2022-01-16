package contexts

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . ContextProvider

type ContextProvider interface {
	GetBuildContext() *BuildContext
	GetDeployContext() *DeployContext
}

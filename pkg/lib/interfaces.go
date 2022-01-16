package lib

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . LibProvider
type LibProvider interface {
	GetDockerLib() Docker
	GetGitLib() Git
}

//counterfeiter:generate . Docker
type Docker interface {
	Build(dockerFilePath string, workingDirectory string, imageName string, imageTag string, additionalBuildArgs []string)
	PullLatest(imageName string)
	Push(imageName string, imageTag string)
	PrintDockerBuild(dockerFilePath string, imageName string, imageTag string, additionalBuildArgs []string) string
	IsImageAvailable(imageName string, imageTag string, useRemoteContainerRegistryCache bool) bool
	RunImage(imageNameAndTag string, envVars map[string]string, additionalArgs map[string]string) error
}

//counterfeiter:generate . Git
type Git interface {
	LoadRepoInformation() *GitState
}

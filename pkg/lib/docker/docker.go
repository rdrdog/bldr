package docker

import (
	"fmt"
	"strings"

	"github.com/Redgwell/bldr/pkg/contexts"
	"github.com/Redgwell/bldr/pkg/lib/process"
	"github.com/sirupsen/logrus"
)

const containerCommitShaLabel = "COMMIT_SHA"
const containerBuildNumber = "BUILD_NUMBER"

type Docker struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *Docker {
	return &Docker{
		logger: logger,
	}
}

func (d *Docker) getBuildArgs(dockerFilePath string, imageName string, imageTag string, buildContext *contexts.BuildContext, additionalBuildArgs []string) string {
	var buildArgs strings.Builder
	buildArgs.WriteString("build ")
	buildArgs.WriteString(fmt.Sprintf("-t %s:%s ", imageName, imageTag))
	buildArgs.WriteString(fmt.Sprintf("-t %s:latest ", imageName))
	buildArgs.WriteString(fmt.Sprintf("--cache-from %s:latest ", imageName))
	buildArgs.WriteString(fmt.Sprintf("--label %s=%s ", containerCommitShaLabel, buildContext.GitContext.FullCommitSha))
	buildArgs.WriteString(fmt.Sprintf("--build-arg %s=%s ", containerCommitShaLabel, buildContext.GitContext.FullCommitSha))
	buildArgs.WriteString(fmt.Sprintf("--build-arg %s=%s ", containerBuildNumber, buildContext.BuildNumber))

	if buildContext.DockerContext.UseBuildKit {
		buildArgs.WriteString("--build-arg BUILDKIT_INLINE_CACHE=1 ")
	}

	// Add any additional build args
	for _, arg := range additionalBuildArgs {
		buildArgs.WriteString(fmt.Sprintf("--build-arg %s ", arg))
	}

	buildArgs.WriteString(fmt.Sprintf("-f %s ", dockerFilePath))
	buildArgs.WriteString(".")

	return buildArgs.String()
}

func (d *Docker) Build(dockerFilePath string, imageName string, imageTag string, buildContext *contexts.BuildContext, additionalBuildArgs []string) {
	buildArgs := d.getBuildArgs(dockerFilePath, imageName, imageTag, buildContext, additionalBuildArgs)
	p := process.
		New("docker", buildContext.PathContext.RepoRootDirectory, d.logger).
		WithArgs(buildArgs)

	if buildContext.DockerContext.UseBuildKit {
		d.logger.Info("Buildkit enabled")
		p.WithEnv("DOCKER_BUILDKIT=1")
	}

	d.logger.Debugf("docker %s", buildArgs)

	_, stdErr, err := p.
		PipeStderrToStdout(). // buildkit sends to stderr....
		Run()

	if err != nil {
		d.logger.Error(stdErr)
		d.logger.Fatalf("Docker build failed: %v", err)
	}
}

func (d *Docker) PullLatest(imageName string) {
	args := fmt.Sprintf("pull %s:latest", imageName)
	_, _, err := process.New("docker", ".", d.logger).WithArgs(args).Run()

	if err == nil {
		d.logger.Infof("Docker pull successful: %s", imageName)
	} else {
		d.logger.Fatalf("Docker pull failed: %v", err)
	}
}

func (d *Docker) Push(imageName string, imageTag string) {
	containerNameWithTag := fmt.Sprintf("%s:%s", imageName, imageTag)
	args := fmt.Sprintf("push %s", containerNameWithTag)
	_, _, err := process.New("docker", ".", d.logger).WithArgs(args).Run()

	if err == nil {
		d.logger.Infof("Docker push successful: %s", containerNameWithTag)
	} else {
		d.logger.Fatalf("Docker push failed: %v", err)
	}
}

func (d *Docker) PrintDockerBuild(dockerFilePath string, imageName string, imageTag string, buildContext *contexts.BuildContext, additionalBuildArgs []string) string {
	return "docker " + d.getBuildArgs(dockerFilePath, imageName, imageTag, buildContext, additionalBuildArgs)
}

func (d *Docker) IsImageAvailable(imageName string, imageTag string, useRemoteContainerRegistryCache bool) bool {
	var output string
	if useRemoteContainerRegistryCache {
		d.logger.Infof("Locating container '%s' for sha '%s' using docker cli", imageName, imageTag)
		args := fmt.Sprintf(
			"manifest inspect %s:%s",
			imageName,
			imageTag,
		)

		var err error
		_, _, err = process.New("docker", ".", d.logger).
			WithArgs(args).
			WithSuppressedOutput().
			Run()

		if err == nil {
			d.logger.Infof("Found container %s:%s", imageName, imageTag)
			return true
		} else {
			d.logger.Infof("Unable to find container %s:%s", imageName, imageTag)
			return false
		}

	} else {
		d.logger.Infof("Locating container '%s' for sha '%s' locally", imageName, imageTag)
		args := fmt.Sprintf("images ls --filter reference=%s*%s --format {{.Tag}}", imageName, imageTag)
		var err error
		output, _, err = process.New("docker", ".", d.logger).
			WithArgs(args).
			WithSuppressedOutput().
			Run()

		if err != nil {
			d.logger.Infof("Unable to find container %s locally", imageName)
			return false
		}

		for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
			if strings.Contains(line, imageTag) {
				d.logger.Infof("Found container %s:%s locally", imageName, imageTag)
				return true
			}
		}

		d.logger.Infof("Unable to find container %s:%s locally", imageName, imageTag)
		return false
	}
}
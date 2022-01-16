package docker

import (
	"fmt"
	"strings"

	"github.com/rdrdog/bldr/internal/services/process"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/sirupsen/logrus"
)

type Docker struct {
	configuration *config.Configuration
	logger        *logrus.Logger
}

func New(configuration *config.Configuration, logger *logrus.Logger) *Docker {
	return &Docker{
		configuration: configuration,
		logger:        logger,
	}
}

func (d *Docker) getBuildArgs(dockerFilePath string, imageName string, imageTag string, additionalBuildArgs []string) string {
	var buildArgs strings.Builder
	buildArgs.WriteString("build ")
	buildArgs.WriteString(fmt.Sprintf("-t %s:%s ", imageName, imageTag))
	buildArgs.WriteString(fmt.Sprintf("-t %s:latest ", imageName))
	buildArgs.WriteString(fmt.Sprintf("--cache-from %s:latest ", imageName))

	if d.configuration.Docker.UseBuildKit {
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

func (d *Docker) Build(dockerFilePath string, workingDirectory string, imageName string, imageTag string, additionalBuildArgs []string) {
	buildArgs := d.getBuildArgs(dockerFilePath, imageName, imageTag, additionalBuildArgs)
	p := process.
		New("docker", workingDirectory, d.logger).
		WithArgs(buildArgs)

	if d.configuration.Docker.UseBuildKit {
		d.logger.Info("🐳 buildkit enabled")
		p.WithEnv("DOCKER_BUILDKIT=1")
	}

	d.logger.Debugf("🐳 docker %s", buildArgs)

	_, stdErr, err := p.
		PipeStderrToStdout(). // buildkit sends to stderr....
		Run()

	if err != nil {
		d.logger.Error(stdErr)
		d.logger.Fatalf("docker build failed: %v", err)
	}
}

func (d *Docker) PullLatest(imageName string) {
	args := fmt.Sprintf("pull %s:latest", imageName)
	_, _, err := process.New("docker", ".", d.logger).WithArgs(args).Run()

	if err == nil {
		d.logger.Infof("🐳 docker pull successful: %s", imageName)
	} else {
		d.logger.Fatalf("docker pull failed: %v", err)
	}
}

func (d *Docker) Push(imageName string, imageTag string) {
	containerNameWithTag := fmt.Sprintf("%s:%s", imageName, imageTag)
	args := fmt.Sprintf("push %s", containerNameWithTag)
	_, _, err := process.New("docker", ".", d.logger).WithArgs(args).Run()

	if err == nil {
		d.logger.Infof("🐳 docker push successful: %s", containerNameWithTag)
	} else {
		d.logger.Fatalf("docker push failed: %v", err)
	}
}

func (d *Docker) PrintDockerBuild(dockerFilePath string, imageName string, imageTag string, additionalBuildArgs []string) string {
	return "docker " + d.getBuildArgs(dockerFilePath, imageName, imageTag, additionalBuildArgs)
}

func (d *Docker) IsImageAvailable(imageName string, imageTag string, useRemoteContainerRegistryCache bool) bool {
	var output string
	if useRemoteContainerRegistryCache {
		d.logger.Infof("locating container '%s' for sha '%s' using docker cli", imageName, imageTag)
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
			d.logger.Infof("found container %s:%s", imageName, imageTag)
			return true
		} else {
			d.logger.Infof("unable to find container %s:%s", imageName, imageTag)
			return false
		}

	} else {
		d.logger.Infof("locating container '%s' for sha '%s' locally", imageName, imageTag)
		args := fmt.Sprintf("images ls --filter reference=%s*%s --format {{.Tag}}", imageName, imageTag)
		var err error
		output, _, err = process.New("docker", ".", d.logger).
			WithArgs(args).
			WithSuppressedOutput().
			Run()

		if err != nil {
			d.logger.Infof("unable to find container %s locally", imageName)
			return false
		}

		for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
			if strings.Contains(line, imageTag) {
				d.logger.Infof("found container %s:%s locally", imageName, imageTag)
				return true
			}
		}

		d.logger.Infof("unable to find container %s:%s locally", imageName, imageTag)
		return false
	}
}

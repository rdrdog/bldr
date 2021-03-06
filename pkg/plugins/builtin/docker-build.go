package builtin

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/extensions"
	"github.com/rdrdog/bldr/pkg/lib"
	"github.com/sirupsen/logrus"
)

const buildArgContainerCommitSha = "COMMIT_SHA"
const buildArgContainerBuildNumber = "BUILD_NUMBER"

type DockerBuild struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	Name          string
	Path          string
	Include       []string
}

func (p *DockerBuild) SetConfig(logger *logrus.Logger, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	return mapstructure.Decode(pluginConfig, p)
}

func (p *DockerBuild) Execute(contextProvider contexts.ContextProvider, extensionsProvider extensions.ExtensionsProvider, libProvider lib.LibProvider) error {
	bc := contextProvider.GetBuildContext()

	imageTag := bc.GitContext.ShortCommitSha
	if p.configuration.Docker.IncludeTimeInImageTag {
		// Go datetime format uses: 01/02 03:04:05PM ‘06 -0700.)
		imageTag += time.Now().UTC().Format("-060102150405")
	}

	imageName := fmt.Sprintf("%s%s", p.configuration.Docker.Registry, p.Name)
	docker := libProvider.GetDockerLib()

	if !p.shouldBuildContainer(bc, docker, imageName, imageTag) {
		p.logger.Infof("🦘 skipping build of target: %s", p.Name)

		return nil
	}

	p.logger.Infof("🧱 building container %s -> %s:%s", p.Name, imageName, imageTag)

	if p.configuration.Docker.UseRemoteContainerRegistryCache {
		docker.PullLatest(imageName)
	}

	buildArgs := []string{
		fmt.Sprintf("%s=\"%s\"", buildArgContainerBuildNumber, bc.BuildNumber),
		fmt.Sprintf("%s=\"%s\"", buildArgContainerCommitSha, bc.GitContext.FullCommitSha),
	}

	docker.Build(p.Path, p.configuration.Paths.RepoRootDirectory, imageName, imageTag, buildArgs)

	if p.configuration.Docker.PushContainers {
		docker.Push(imageName, imageTag)
		docker.Push(imageName, "latest")
	} else {
		p.logger.Infof("⏭  skipping container push for %s", imageName)
	}

	bc.ArtefactManifest.AddArtefact(p.Name, fmt.Sprintf("%s:%s", imageName, imageTag))

	return nil
}

func (p *DockerBuild) isAffectedByDiff(diffFilePaths []string) bool {
	for _, pathGlob := range p.Include {

		for _, diffFile := range diffFilePaths {
			isMatch, _ := filepath.Match(pathGlob, diffFile)

			if isMatch {
				p.logger.Debugf("found match on target '%s' (glob %s matched on file %s)", p.Name, pathGlob, diffFile)
				return true
			}
		}
	}

	return false
}

func (p *DockerBuild) shouldBuildContainer(bc *contexts.BuildContext, docker lib.Docker, imageName string, imageTag string) bool {
	if !bc.GitContext.CanDetectChanges() {
		p.logger.Info("git context not in a state to detect changes - build is required")
		return true
	}

	if p.isAffectedByDiff(bc.GitContext.ChangesSinceMainBranch) {
		p.logger.Infof("🪢  target %s found to be affected by changes", p.Name)
		return true
	} else if docker.IsImageAvailable(imageName, bc.GitContext.MainBranchForkPointShort, p.configuration.Docker.UseRemoteContainerRegistryCache) {
		p.logger.Infof("🔎 found container for branch base fork commit: %s - no need to build it", bc.GitContext.MainBranchForkPointShort)
		bc.ArtefactManifest.AddArtefact(p.Name, fmt.Sprintf("%s:%s", imageName, bc.GitContext.MainBranchForkPointShort))
		return false
	} else {
		p.logger.Infof("🥷 could not locate image for branch base fork commit: %s - building image", bc.GitContext.MainBranchForkPointShort)
		return true
	}
}

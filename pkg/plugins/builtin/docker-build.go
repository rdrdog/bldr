package builtin

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rdrdog/bldr/pkg/config"
	"github.com/rdrdog/bldr/pkg/contexts"
	"github.com/rdrdog/bldr/pkg/lib/docker"
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

func (p *DockerBuild) SetConfig(logger *logrus.Logger, targetName string, configuration *config.Configuration, pluginConfig map[string]interface{}) error {
	p.configuration = configuration
	p.logger = logger
	p.Name = targetName
	return mapstructure.Decode(pluginConfig, p)
}

func (p *DockerBuild) Execute(contextProvider *contexts.ContextProvider) error {
	p.logger.Infof("Running docker build with config: Path: %s, Include: %v", p.Path, p.Include)

	bc := contextProvider.BuildContext

	imageTag := bc.GitContext.ShortCommitSha
	if p.configuration.Docker.IncludeTimeInImageTag {
		// Go datetime format uses: 01/02 03:04:05PM ‘06 -0700.)
		imageTag += time.Now().UTC().Format("-060102150405")
	}

	imageName := fmt.Sprintf("%s%s", p.configuration.Docker.Registry, p.Name)
	docker := docker.New(p.configuration, p.logger)

	shouldBuildContainer := true
	/*
	   if bc.GitContext.CanDetectChanges() {
	     if (cm.IsAffectedByDiff(gitState.ChangesSinceMainBranch))
	     {
	         Log.Information($"Container {cm.Name} found to be affected by changes");
	         shouldBuildContainer = true;
	     }
	     else if (_docker.IsImageAvailable(container.GetFQImageName(), gitState.MainBranchForkPoint, cfg.QueryRemoteRegistryForTags))
	     {
	         Log.Information($"Found container for branch base fork commit: ${gitState.MainBranchForkPoint} - no need to build it");
	         container.SetImageTagFromCommitSha(gitState.MainBranchForkPoint, false);
	         manifest.WithImage(container);
	         shouldBuildContainer = false;
	     }
	     else
	     {
	         Log.Information($"Could not locate image for branch base fork commit: {gitState.MainBranchForkPoint}. Building image.");
	         shouldBuildContainer = true;
	     }
	   }
	*/

	if shouldBuildContainer {

		p.logger.Infof("🧱 Building container %s -> %s:%s", p.Name, imageName, imageTag)

		if p.configuration.Docker.UseRemoteContainerRegistryCache {
			docker.PullLatest(imageName)
		}

		buildArgs := []string{
			fmt.Sprintf("%s=\"%s\"", buildArgContainerBuildNumber, bc.BuildNumber),
			fmt.Sprintf("%s=\"%s\"", buildArgContainerCommitSha, bc.GitContext.FullCommitSha),
		}

		docker.Build(p.Path, bc.PathContext.RepoRootDirectory, imageName, imageTag, buildArgs)

		if p.configuration.Docker.PushContainers {
			docker.Push(imageName, imageTag)
			docker.Push(imageName, "latest")
		} else {
			p.logger.Infof("⏭  skipping container push for %s", imageName)
		}

		bc.ArtefactManifest.AddArtefact(p.Name, fmt.Sprintf("%s:%s", imageName, imageTag))
	} else {
		p.logger.Infof("🦘 Skipping build of target: %s", p.Name)
	}

	return nil
}

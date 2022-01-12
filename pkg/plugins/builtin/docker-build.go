package builtin

import (
	"fmt"
	"time"

	"github.com/Redgwell/bldr/pkg/contexts"
	"github.com/Redgwell/bldr/pkg/lib/docker"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type DockerBuild struct {
	logger  *logrus.Logger
	Name    string
	Path    string
	Include []string
}

func (p *DockerBuild) SetConfig(logger *logrus.Logger, targetName string, pluginConfig map[string]interface{}) error {
	p.logger = logger
	p.Name = targetName
	return mapstructure.Decode(pluginConfig, p)
}

func (p *DockerBuild) Execute(contextProvider *contexts.ContextProvider) error {
	p.logger.Infof("Running docker build with config: Path: %s, Include: %v", p.Path, p.Include)

	p.logger.Infof("context: %v", contextProvider.BuildContext.BuildNumber)

	bc := contextProvider.BuildContext

	imageTag := bc.GitContext.ShortCommitSha
	if bc.DockerContext.IncludeTimeInTags {
		// Go datetime format uses: 01/02 03:04:05PM ‘06 -0700.)
		imageTag += time.Now().UTC().Format("-060201150405")
	}

	imageName := fmt.Sprintf("%s%s", bc.DockerContext.Registry, p.Name)
	docker := docker.New(p.logger)

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

		if bc.DockerContext.UseRemoteContainerRegistryCache {
			docker.PullLatest(imageName)
		}

		docker.Build(p.Path, imageName, imageTag, bc, nil)

		if bc.DockerContext.PushContainers {
			docker.Push(imageName, imageTag)
			docker.Push(imageName, "latest")
		} else {
			p.logger.Infof("⏭  skipping container push for %s", imageName)
		}

		contextProvider.BuildContext.ArtefactManifest.AddArtefact(p.Name, imageName)
	} else {
		p.logger.Infof("🦘 Skipping build of target: %s", p.Name)
	}

	return nil
}

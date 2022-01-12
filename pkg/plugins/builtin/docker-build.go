package builtin

import (
	"github.com/Redgwell/bldr/pkg/contexts"
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

func (p *DockerBuild) Execute(contextProvider *contexts.ContextProvider /*projectName string, targetName string*/) error {
	p.logger.Infof("Running docker build with config: Path: %s, Include: %v", p.Path, p.Include)

	shouldBuildContainer := true
	/*
		if (gitState.CanDetectChanges()) {
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
		// Log.Information($"⚒️ Building container {container.Name} -> {container.GetFQImageNameAndTag()}");
		// manifest.WithImage(container).Persist(_fileSystem, manifestFilePath);

		// if (cfg.UseRemoteContainerRegistryCache)
		// {
		// 		_docker.DockerPull(container.GetFQImageNameAndLatestTag(), true);
		// }
		// _docker.DockerBuild(PathProvider.RootWorkingDirectory, cfg, cm, container, gitState.CommitSha);

		// PushContainer(container, cfg);

	} else {
		p.logger.Infof("🦘 Skipping build of target: %s", p.Name)

	}
	return nil
}

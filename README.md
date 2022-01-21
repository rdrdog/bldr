# bldr
[![Build](https://github.com/rdrdog/bldr/actions/workflows/build.yaml/badge.svg)](https://github.com/rdrdog/bldr/actions/workflows/build.yaml)


**Bldr** is a framework to make building containers in a mono-repo easy!

**Dplyr** is a companion framework to make deploying these containers to Kubernetes (or other cloud services) easy too!

## Standard use cases
You should consider using this framework if:
- You're running a platform team in an organisation with multiple squads
- You want to adopt a standard approach to build and deploy pipelines
- Your squads don't yet have the knowledge/maturity/desire to maintain their own helm charts (for Kubernetes)
- You want to continuously improve your build and deployment pipelines in a way that can be rolled out across all squads
- You want a pluggable architecture for your deployment pipelines that will support growth as you target new cloud services

## Getting started
1. Download the latest bldr binary from the releases page.
1. Install the bldr binary somewhere on your path.
1. Create a pipeline-config.yaml file in your repository root. Use the [sample](samples/pipeline-config-simple.yaml) to get started.
1. Run `bldr build` to run your build.
    - This will generate a default bldr.yaml file in the root of your repository with safe defaults. Before committing, you will need to set your container registry URL for CI builds
1. Run `bldr deploy` to verify your deployment succeeds

You're now running your build and deploy pipeline on your local machine!
Next, you will need to get your GitHub Actions pipeline working:

1. Create the .github/workflows folders
    ```
    mkdir -p .github/workflows
    ```
1. Copy the [sample github actions yaml](samples/example-gh-build-and-deploy.yml) to your .github/workflows folder
1. Modify the sample to match your cloud environment names

> TODO - specify how clusters are configured, bldr.yaml configuration value details, and pipeline-config.yaml details

## Motivation
Modern engineering squads are usually familiar with concepts and tooling like infra-as-code, containers, Kubernetes, test automation, mono-repos, secure secret management etc.

However many organisations are not currently practicing these modern engineering concepts in all of their squads. Often, squads may have heard of some of these practices, but they may not necessarily understand them well enough to be able to adopt them with confidence.

Platform teams can help to bridge this knowledge gap by providing a platform in which squads can develop on, with safe defaults, easy build and deployment pipeline setup, and easy bootstrapping of default infrastructure.

This tooling helps support platform teams to create a standard pattern to build and deploy multiple containers from mono-repos into targets such as Kubernetes, without the squads having to be well versed in the intricacies of creating their own helm charts, managing secrets securely, or any of the other complications.

## Concepts
There are a few concepts to understand as you dive into this tooling.

### 1. Plugins
Plugins are implementations of specific types of steps in your build or deployment pipeline. Checkout the [list of built in plugins](https://github.com/rdrdog/bldr/tree/main/pkg/plugins/builtin).

### 2. Extensions
Extensions provide cross cutting functionality that can be used by any plugin as the pipelines run. Checkot the [list of built in extensions](https://github.com/rdrdog/bldr/tree/main/pkg/extensions/builtin)
> Right now, the only extension is for secret loading, but there are plans to add more soon!

### 3. Monorepos
This tooling works great for repositories with single artifacts, but really, it shines for monorepo scenarios.

One of the common complaints with monorepos is that builds and deployments take too long, since even the smallest code changes require everything to be built and deployed.

This tooling provides optimisations to solve this problem by just building the containers that have changed on the branch (as opposed to changes for each commit). At deploy time, all containers in the repository are deployed. However, this means that for containers that _aren't_ built on a branch, we deploy the previously built container from the main branch.

## FAQ
_Why not just use Github actions? (or Azure DevOps, GitlabCI, TravicsCI, etc.)_

While most modern CI tooling provides specific actions/jobs to do things like building containers, or deploying to Kubernetes, they:
- lack optimisations around building just the containers that have changed in a mono-repo scenario
- make deploying to Kubernetes using an organisation managed helm chart quite complex
- lack the ability to build and deploy locally using the same build/deploy pipelines as used in the cloud (fast feedback!!!)
- lack standardised secure secret management for both local and cloud environments
- end up being a mishmash of squad specific YAML/pipelines, with little to now consistency making it hard to add future optimisations that all squads can benefit from

_I have already built this stuff before using Azure DevOps (or Github Actions, GitlabCI, TravicsCI, etc.)!_
- Great! Keep running with that, supporting your squads :) This tooling is here to give you a head start if you don't want to create this all from scratch

## Development

Running tests:
```
make test
```

Generating fakes (for new, or changed interfaces):
```
make generate
```

Running bldr:

```
go run cmd/bldr/main.go
```

Running dplyr:

```
go run cmd/dplyr/main.go
```

## Plugins

- Some default initial plugins:
  - deploy
    - DockerRun
    - K8sDeploy


- Could look to use the go plugin system: https://medium.com/learning-the-go-programming-language/writing-modular-go-programs-with-plugins-ec46381ee1a9
  - Initialise plugins by building them using `go build --buildmode=plugin -o /plugins/something.so github.com/something/plugin.go`

- Detecting diffs should be a plugin, so that we can be smarter than just relying on devs to populate `include` properly


- Deployment:
  - How do we handle helm charts for deployment?
    - Platform teams will want their own helm chart
      - Could be a platform-team specific repo
      - manifest.yaml output can be used to drive it
      - GH action like helm-deploy could be used to deploy the helm chart from the manifest:
        https://github.com/marketplace/actions/helm-deploy


TODO:
- k8s deploy using helm chart from another repo
- missing tests
- docs
- example usage repo
- standardise usage of .Fatal and returning errors



GH DA
- We could split the logic into separate individual custom composite actions:
  Build:
    - BuildPathContextLoader:
      - not needed, since not needed, since GH can publish artfefacts arbitrarily

    - GitContextLoader
98      - use step state to store outputs (https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions)
      - run each of the git commands in GH actions

    - Docker-build
      - Go code uses pipeline-config to drive it
      - GH actions would need to have a job to build each container
      - There's no easy way in GH actions to run jobs only if files have changed (but we could create a reusable workflow that does this based on the GitContextLoader output - it could take the name, path, and include glob(s) just like we have in pipeline-config)
      - Tag (local vs CI logic for time in image tag)
      - push (only on CI)
      - registry as part of name (local vs CI config)
      - using remote container cache (local vs CI config)

  Deploy:
    - LoadDeployContext:
      - EnvironmentName is a param to the WF
      - Manifest can be read from artefacts, and loaded into step state

    - Run container:
      - we can get secrets using key/values as inputs, which we load from AWS secret manager, and load as env vars
      - can then get the image tag from the manifest step (image key as input)

    - Deploy to k8s
      - we can get secrets using key/values as inputs, which we load from AWS secret manager, and load as env vars
      - we can clone another repo's helm chart (this repo can be an input param)
      - the cluster config for a given env could be mapped against gh environments/variables, or even in the action has hard-coded mappings/config


- Problems:
  - Reusuable workflows are not yet supported (https://github.com/nektos/act/issues/826)
    - Perhaps this won't be a problem if composite workflows are supported?

  - Output sharing
    - We can have each thing publish data to a 'manifest' folder, where files can be keys, and contents can be values
    - This way, we can share state between jobs without requiring specific job identification
      - Not sure how concurrency works?

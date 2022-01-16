# bldr
[![Build](https://github.com/rdrdog/bldr/actions/workflows/build.yaml/badge.svg)](https://github.com/rdrdog/bldr/actions/workflows/build.yaml)


## Development

Running tests:
```
go test ./... -v
```

Generating fakes (for new, or changed interfaces):
```
go generate ./...
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

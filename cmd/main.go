package main

import (
	"fmt"

	"github.com/Redgwell/bldr/internal/pipeline"
)

func main() {
	fmt.Println("Hello, World!")

	pp, _ := pipeline.NewPluginPipeline() // todo - pass in config
	pp.Run()
}

/*
TODO:
- main entrypoint to load config
	- Use env override stuff from WP
	- load config from .bldr
	- use default plugin suite for the command ('build', or 'deploy')
- Init the pluginPipeline from here then execute it
*/

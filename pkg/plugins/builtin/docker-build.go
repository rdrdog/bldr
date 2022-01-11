package builtin

import "fmt"

type DockerBuild struct {
}

func (p *DockerBuild) Execute() error {
	fmt.Println("Running docker build")
	return nil
}

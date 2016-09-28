package images

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy/shared/containers"
)

// Run the images flow
func Run(c *cli.Context) {
	for _, container := range containers.AllContainers {
		fmt.Println(container.Image)
	}
}

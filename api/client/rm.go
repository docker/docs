package client

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	Cli "github.com/docker/docker/cli"
	flag "github.com/docker/docker/pkg/mflag"
)

// CmdRm removes one or more containers.
//
// Usage: docker rm [OPTIONS] CONTAINER [CONTAINER...]
func (cli *DockerCli) CmdRm(args ...string) error {
	cmd := Cli.Subcmd("rm", []string{"CONTAINER [CONTAINER...]"}, Cli.DockerCommands["rm"].Description, true)
	v := cmd.Bool([]string{"v", "-volumes"}, false, "Remove the volumes associated with the container")
	link := cmd.Bool([]string{"l", "-link"}, false, "Remove the specified link")
	force := cmd.Bool([]string{"f", "-force"}, false, "Force the removal of a running container (uses SIGKILL)")
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args, true)

	var errs []string
	for _, name := range cmd.Args() {
		if name == "" {
			return fmt.Errorf("Container name cannot be empty")
		}
		name = strings.Trim(name, "/")

		options := types.ContainerRemoveOptions{
			ContainerID:   name,
			RemoveVolumes: *v,
			RemoveLinks:   *link,
			Force:         *force,
		}

		if err := cli.client.ContainerRemove(options); err != nil {
			errs = append(errs, fmt.Sprintf("Failed to remove container (%s): %s", name, err))
		} else {
			fmt.Fprintf(cli.out, "%s\n", name)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "\n"))
	}
	return nil
}

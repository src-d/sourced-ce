package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
)

type logsCmd struct {
	Command `name:"logs" short-description:"Fetch the logs of source{d} components" long-description:"Fetch the logs of source{d} components"`

	Follow bool `short:"f" long:"follow" description:"Follow log output"`
	Args   struct {
		Components []string `positional-arg-name:"component" description:"Component names from where to fetch logs"`
	} `positional-args:"yes"`
}

func (c *logsCmd) Execute(args []string) error {
	command := []string{"logs"}

	if c.Follow {
		command = append(command, "--follow")
	}

	if components := c.Args.Components; len(components) > 0 {
		command = append(command, components...)
	}

	return compose.Run(context.Background(), command...)
}

func init() {
	rootCmd.AddCommand(&logsCmd{})
}

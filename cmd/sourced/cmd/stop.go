package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
)

type stopCmd struct {
	Command `name:"stop" short-description:"Stop running containers" long-description:"Stop running containers without removing them.\nThey can be started again with 'start'."`
}

func (c *stopCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "stop")
}

func init() {
	rootCmd.AddCommand(&stopCmd{})
}

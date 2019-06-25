package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
)

type stopCmd struct {
	Command `name:"stop" short-description:"Stop any running components" long-description:"Stop any running components without removing them.\nThey can be started again with 'start'."`
}

func (c *stopCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "stop")
}

func init() {
	rootCmd.AddCommand(&stopCmd{})
}

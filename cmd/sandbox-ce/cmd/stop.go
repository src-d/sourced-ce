package cmd

import (
	"context"

	"github.com/smacker/superset-compose/cmd/sandbox-ce/compose"
)

type stopCmd struct {
	Command `name:"stop" short-description:"Stop"`
}

func (c *stopCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "stop")
}

func init() {
	rootCmd.AddCommand(&stopCmd{})
}

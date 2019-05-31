package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sandbox-ce/compose"
)

type statusCmd struct {
	Command `name:"status" short-description:"Shows status of the components" long-description:"Shows status of the components"`
}

func (c *statusCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "ps")
}

func init() {
	rootCmd.AddCommand(&statusCmd{})
}

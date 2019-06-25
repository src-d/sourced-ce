package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
)

type statusCmd struct {
	Command `name:"status" short-description:"Show the status of all components" long-description:"Show the status of all components"`
}

func (c *statusCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "ps")
}

func init() {
	rootCmd.AddCommand(&statusCmd{})
}

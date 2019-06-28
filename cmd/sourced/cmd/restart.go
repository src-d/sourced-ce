package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
)

type restartCmd struct {
	Command `name:"restart" short-description:"Update current installation according to the active docker compose file" long-description:"Update current installation according to the active docker compose file."`
}

func (c *restartCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "up", "--force-recreate")
}

func init() {
	rootCmd.AddCommand(&restartCmd{})
}

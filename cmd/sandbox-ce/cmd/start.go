package cmd

import (
	"context"

	"github.com/src-d/superset-compose/cmd/sandbox-ce/compose"
)

type startCmd struct {
	Command `name:"start" short-description:"Start stopped containers" long-description:"Start stopped containers.\nThe containers must be initialized before with 'install'."`
}

func (c *startCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "start")
}

func init() {
	rootCmd.AddCommand(&startCmd{})
}

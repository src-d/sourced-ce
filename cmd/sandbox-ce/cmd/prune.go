package cmd

import (
	"context"

	"github.com/src-d/superset-compose/cmd/sandbox-ce/compose"
)

type pruneCmd struct {
	Command `name:"prune" short-description:"Stop and remove containers and resources" long-description:"Stops containers and removes containers, networks, and volumes created by 'install'.\nImages are not deleted from the system."`
}

func (c *pruneCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "down", "--volumes")
}

func init() {
	rootCmd.AddCommand(&pruneCmd{})
}

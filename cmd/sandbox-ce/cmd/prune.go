package cmd

import (
	"context"

	"github.com/smacker/superset-compose/cmd/sandbox-ce/compose"
)

type pruneCmd struct {
	Command `name:"prune" short-description:"Prune"`
}

func (c *pruneCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "down")
}

func init() {
	rootCmd.AddCommand(&pruneCmd{})
}

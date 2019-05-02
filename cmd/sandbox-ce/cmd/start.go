package cmd

import (
	"context"

	"github.com/smacker/superset-compose/cmd/sandbox-ce/compose"
)

type startCmd struct {
	Command `name:"start" short-description:"Start"`
}

func (c *startCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "start")
}

func init() {
	rootCmd.AddCommand(&startCmd{})
}

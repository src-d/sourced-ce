package cmd

import (
	"context"

	"github.com/src-d/superset-compose/cmd/sandbox-ce/compose"
)

type installCmd struct {
	Command `name:"install" short-description:"Install"`
}

func (c *installCmd) Execute(args []string) error {
	err := compose.Run(context.Background(),
		"run", "--rm", "superset", "./docker-init.sh")

	if err != nil {
		return err
	}

	err = compose.Run(context.Background(), "up")
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(&installCmd{})
}

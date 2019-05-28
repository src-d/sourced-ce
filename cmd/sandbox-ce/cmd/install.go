package cmd

import (
	"context"
	"time"

	"github.com/src-d/superset-compose/cmd/sandbox-ce/compose"
)

type installCmd struct {
	Command `name:"install" short-description:"Install and initialize containers" long-description:"Install, initialize, and start all the required docker containers, networks, volumes, and images."`
}

func (c *installCmd) Execute(args []string) error {
	if err := compose.Run(context.Background(), "up", "--detach"); err != nil {
		return err
	}

	return OpenUI(time.Minute)
}

func init() {
	rootCmd.AddCommand(&installCmd{})
}

package cmd

import (
	"context"
	"time"

	"github.com/src-d/superset-compose/cmd/sandbox-ce/compose"
)

type startCmd struct {
	Command `name:"start" short-description:"Start stopped containers" long-description:"Start stopped containers.\nThe containers must be initialized before with 'install'."`
}

func (c *startCmd) Execute(args []string) error {
	done := OpenUI(10 * time.Second)

	if err := compose.Run(context.Background(), "start"); err != nil {
		return err
	}

	if err := <-done; err != nil {
		return err
	}

	return nil

}

func init() {
	rootCmd.AddCommand(&startCmd{})
}

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
	if err := compose.Run(context.Background(), "start"); err != nil {
		return err
	}

	return OpenUI(30 * time.Second)

}

func init() {
	rootCmd.AddCommand(&startCmd{})
}

package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
)

type sqlCmd struct {
	Command `name:"sql" short-description:"Open a MySQL client connected to gitbase" long-description:"Open a MySQL client connected to gitbase"`
}

func (c *sqlCmd) Execute(args []string) error {
	return compose.Run(context.Background(), "exec", "gitbase", "mysql")
}

func init() {
	rootCmd.AddCommand(&sqlCmd{})
}

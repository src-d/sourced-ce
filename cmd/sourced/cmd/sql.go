package cmd

import (
	"context"
	"os"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"

	"golang.org/x/crypto/ssh/terminal"
)

type sqlCmd struct {
	Command `name:"sql" short-description:"Open a MySQL client connected to a SQL interface for Git" long-description:"Open a MySQL client connected to a SQL interface for Git"`

	Args struct {
		Query string `positional-arg-name:"query" description:"SQL query to be run by the SQL interface for Git"`
	} `positional-args:"yes"`
}

func (c *sqlCmd) Execute(args []string) error {
	command := []string{"exec"}
	if !terminal.IsTerminal(int(os.Stdout.Fd())) || !terminal.IsTerminal(int(os.Stdin.Fd())) {
		command = append(command, "-T")
	}
	command = append(command, "gitbase", "mysql")
	if c.Args.Query != "" {
		command = append(command, "--execute", c.Args.Query)
	}

	return compose.Run(context.Background(), command...)
}

func init() {
	rootCmd.AddCommand(&sqlCmd{})
}

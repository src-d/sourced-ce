package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gopkg.in/src-d/go-cli.v0"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type orgsCmd struct {
	cli.PlainCommand `name:"orgs" short-description:"Manages services to analyze code from GitHub organizations"`
}

type orgsInitCmd struct {
	Command `name:"init" short-description:"Install and initialize containers to analyze github organizations" long-description:"Install, initialize, and start all the required docker containers, networks, volumes, and images.\n\nThe orgs argument must a list of the organizations to be analyzed."`

	Token string `short:"t" long:"token" description:"Github token for the passed organizations. It should be granted with 'repo' and 'read:org' scopes." required:"true"`
	Args  struct {
		Orgs []string `required:"yes"`
	} `positional-args:"yes" required:"1"`
}

func (c *orgsInitCmd) Execute(args []string) error {
	dir, err := workdir.InitWithOrgs(c.Args.Orgs, c.Token)
	if err != nil {
		return err
	}

	// Before setting a new workdir, stop the current containers
	compose.Run(context.Background(), "stop")

	err = workdir.SetActive(dir)
	if err != nil {
		return err
	}

	fmt.Printf("docker-compose working directory set to %s\n", strings.Join(c.Args.Orgs, ","))

	if err := compose.Run(context.Background(), "up", "--detach"); err != nil {
		return err
	}

	return OpenUI(60 * time.Minute)
}

func init() {
	c := rootCmd.AddCommand(&orgsCmd{})
	c.AddCommand(&orgsInitCmd{})
}

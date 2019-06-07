package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type pruneCmd struct {
	Command `name:"prune" short-description:"Stop and remove containers and resources" long-description:"Stops containers and removes containers, networks, and volumes created by 'init' for the current working directory.\nTo delete resources for all working directories pass --all flag.\nImages are not deleted unless you specify the --images flag."`

	All    bool `short:"a" long:"all" description:"Remove containers and resources for all working directories"`
	Images bool `long:"images" description:"Remove docker images"`
}

func (c *pruneCmd) Execute(args []string) error {
	if !c.All {
		return c.pruneActive()
	}

	dirs, err := workdir.ListRepoDirs()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if err := workdir.SetActive(dir); err != nil {
			return err
		}

		if err = c.pruneActive(); err != nil {
			return err
		}
	}

	return nil
}

func (c *pruneCmd) pruneActive() error {
	a := []string{"down", "--volumes"}
	if c.Images {
		a = append(a, "--rmi", "all")
	}

	return compose.Run(context.Background(), a...)
}

func init() {
	rootCmd.AddCommand(&pruneCmd{})
}

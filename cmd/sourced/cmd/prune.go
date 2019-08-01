package cmd

import (
	"context"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type pruneCmd struct {
	Command `name:"prune" short-description:"Stop and remove components and resources" long-description:"Stops containers and removes containers, networks, volumes and configuration created by 'init' for the current working directory.\nTo delete resources for all working directories pass --all flag.\nImages are not deleted unless you specify the --images flag."`

	All    bool `short:"a" long:"all" description:"Remove containers and resources for all working directories"`
	Images bool `long:"images" description:"Remove docker images"`
}

func (c *pruneCmd) Execute(args []string) error {
	workdirHandler, err := workdir.NewHandler()
	if err != nil {
		return err
	}

	if !c.All {
		return c.pruneActive(workdirHandler)
	}

	wds, err := workdirHandler.List()
	if err != nil {
		return err
	}

	for _, wd := range wds {
		if err := workdirHandler.SetActive(wd); err != nil {
			return err
		}

		if err = c.pruneActive(workdirHandler); err != nil {
			return err
		}
	}

	return nil
}

func (c *pruneCmd) pruneActive(workdirHandler *workdir.Handler) error {
	a := []string{"down", "--volumes"}
	if c.Images {
		a = append(a, "--rmi", "all")
	}

	if err := compose.Run(context.Background(), a...); err != nil {
		return err
	}

	wd, err := workdirHandler.Active()
	if err != nil {
		return err
	}

	if err := workdirHandler.Remove(wd); err != nil {
		return err
	}

	return workdirHandler.UnsetActive()
}

func init() {
	rootCmd.AddCommand(&pruneCmd{})
}

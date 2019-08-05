package cmd

import (
	"fmt"

	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type workdirsCmd struct {
	Command `name:"workdirs" short-description:"List all working directories" long-description:"List all the previously initialized working directories."`
}

func (c *workdirsCmd) Execute(args []string) error {
	workdirHandler, err := workdir.NewHandler()
	if err != nil {
		return err
	}

	wds, err := workdirHandler.List()
	if err != nil {
		return err
	}

	active, err := workdirHandler.Active()
	if err != nil {
		return err
	}

	for _, wd := range wds {
		if wd.Path == active.Path {
			fmt.Printf("* %s\n", wd.Name)
		} else {
			fmt.Printf("  %s\n", wd.Name)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(&workdirsCmd{})
}

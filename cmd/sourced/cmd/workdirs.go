package cmd

import (
	"fmt"

	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type workdirsCmd struct {
	Command `name:"workdirs" short-description:"List all working directories" long-description:"List all the previously initialized working directories."`
}

func (c *workdirsCmd) Execute(args []string) error {
	dirs, err := workdir.List()
	if err != nil {
		return err
	}

	active, err := workdir.Active()
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir == active {
			fmt.Printf("* %s\n", dir)
		} else {
			fmt.Printf("  %s\n", dir)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(&workdirsCmd{})
}

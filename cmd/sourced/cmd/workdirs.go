package cmd

import (
	"fmt"
	"os"

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

	// active directory not necessary exists
	var activePath = ""
	active, err := workdirHandler.Active()
	if err == nil {
		activePath = active.Path
	} else if !isNotExist(err) {
		return err
	}

	for _, wd := range wds {
		if wd.Path == activePath {
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

func isNotExist(err error) bool {
	if os.IsNotExist(err) {
		return true
	}

	if cause, ok := err.(causer); ok {
		return isNotExist(cause.Cause())
	}

	return false
}

type causer interface {
	Cause() error
}

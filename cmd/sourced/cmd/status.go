package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type statusCmd struct {
	Command `name:"status" short-description:"Show the list of working directories and the current deployment" long-description:"Show the list of working directories and the current deployment"`
}

type statusAllCmd struct {
	Command `name:"all" short-description:"Show all the available status information" long-description:"Show all the available status information"`
}

func (c *statusAllCmd) Execute(args []string) error {
	fmt.Print("List of all working directories:\n")

	err := printWorkdirsCmd()
	if err != nil {
		return err
	}

	active, err := activeWorkdir()
	if isNotExist(err) {
		// skip printing the config and components when there is no active dir
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Print("\nConfiguration used for the active working directory:\n\n")

	err = printConfigCmd(active)
	if err != nil {
		return err
	}

	fmt.Print("\nStatus of all components:\n\n")

	err = printComponentsCmd()
	if err != nil {
		return err
	}

	return nil
}

type statusComponentsCmd struct {
	Command `name:"components" short-description:"Show the status of the components containers" long-description:"Show the status of the components containers"`
}

func (c *statusComponentsCmd) Execute(args []string) error {
	return printComponentsCmd()
}

func printComponentsCmd() error {
	return compose.Run(context.Background(), "ps")
}

type statusWorkdirsCmd struct {
	Command `name:"workdirs" short-description:"List all working directories" long-description:"List all the previously initialized working directories"`
}

func (c *statusWorkdirsCmd) Execute(args []string) error {
	return printWorkdirsCmd()
}

func printWorkdirsCmd() error {
	workdirHandler, err := workdir.NewHandler()
	if err != nil {
		return err
	}

	wds, err := workdirHandler.List()
	if err != nil {
		return err
	}

	activePath, err := activeWorkdir()
	// active directory does not necessarily exist
	if err != nil && !isNotExist(err) {
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

type statusConfigCmd struct {
	Command `name:"config" short-description:"Show the configuration for the active working directory" long-description:"Show the docker-compose environment variables configuration for the active working directory"`
}

func (c *statusConfigCmd) Execute(args []string) error {
	active, err := activeWorkdir()
	if err != nil {
		return err
	}

	return printConfigCmd(active)
}

func printConfigCmd(path string) error {
	content, err := ioutil.ReadFile(filepath.Join(path, ".env"))
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", content)

	return nil
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

func activeWorkdir() (string, error) {
	workdirHandler, err := workdir.NewHandler()
	if err != nil {
		return "", err
	}

	active, err := workdirHandler.Active()
	if err != nil {
		return "", err
	}

	return active.Path, err
}

func init() {
	c := rootCmd.AddCommand(&statusCmd{})

	c.AddCommand(&statusAllCmd{})
	c.AddCommand(&statusComponentsCmd{})
	c.AddCommand(&statusWorkdirsCmd{})
	c.AddCommand(&statusConfigCmd{})
}

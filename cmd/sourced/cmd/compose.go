package cmd

import (
	"fmt"
	"strconv"

	composefile "github.com/src-d/sourced-ce/cmd/sourced/compose/file"

	"gopkg.in/src-d/go-cli.v0"
)

type composeCmd struct {
	cli.PlainCommand `name:"compose" short-description:"Manage source{d} docker compose files" long-description:"Manage source{d} docker compose files"`
}

type composeDownloadCmd struct {
	Command `name:"download" short-description:"Download docker compose files" long-description:"Download docker compose files. By default the command downloads the file for this binary version.\n\nUse the 'version' argument to choose a specific revision from\nthe https://github.com/src-d/sourced-ce repository, or to set a\nURL to a docker-compose.yml file.\n\nExamples:\n\nsourced compose download\nsourced compose download v0.0.1\nsourced compose download master\nsourced compose download https://raw.githubusercontent.com/src-d/sourced-ce/master/docker-compose.yml"`

	Args struct {
		Version string `positional-arg-name:"version" description:"Either a revision (tag, full sha1) or a URL to a docker-compose.yml file"`
	} `positional-args:"yes"`
}

func (c *composeDownloadCmd) Execute(args []string) error {
	v := c.Args.Version
	if v == "" {
		v = version
	}

	err := composefile.Download(v)
	if err != nil {
		return err
	}

	fmt.Println("Docker compose file successfully downloaded to your ~/.sourced/compose-files directory. It is now the active compose file.")
	fmt.Println("To update your current installation use `sourced restart`")
	return nil
}

type composeListCmd struct {
	Command `name:"list" short-description:"List the downloaded docker compose files" long-description:"List the downloaded docker compose files"`
}

func (c *composeListCmd) Execute(args []string) error {
	active, err := composefile.Active()
	if err != nil {
		return err
	}

	files, err := composefile.List()
	if err != nil {
		return err
	}

	for index, file := range files {
		fmt.Printf("[%d]", index)
		if file == active {
			fmt.Printf("* %s\n", file)
		} else {
			fmt.Printf("  %s\n", file)
		}
	}

	return nil
}

type composeSetDefaultCmd struct {
	Command `name:"set" short-description:"Set the active docker compose file" long-description:"Set the active docker compose file"`

	Args struct {
<<<<<<< HEAD
=======
		// Version string `positional-arg-name:"version" description:"Either a revision (tag, full sha1) or a URL to a docker-compose.yml file"`
>>>>>>> f8b01f1... Added index numbers for docker compose files. Signed-off-by: Cihan Mete Bahadir <c.mete.bahadir@gmail.com>
		Index string `positional-arg-name:"index" description:"Index of the docker compose file returned from 'sourced compose list'"`
	} `positional-args:"yes" required:"yes"`
}

func (c *composeSetDefaultCmd) Execute(args []string) error {
	files, err := composefile.List()

	if err != nil {
		return err
	}

	index, err := strconv.Atoi(c.Args.Index)

	if err == nil && index >= 0 && index < len(files) {
		active := files[index]
		err = composefile.SetActive(active)
		fmt.Println("Active docker compose file was changed successfully.")
		fmt.Println("To update your current installation use `sourced restart`")
	} else {
		fmt.Println("Provide the index of the docker compose file in 'sourced compose list'")
	}

	return nil
}

func init() {
	c := rootCmd.AddCommand(&composeCmd{})
	c.AddCommand(&composeDownloadCmd{})
	c.AddCommand(&composeListCmd{})
	c.AddCommand(&composeSetDefaultCmd{})
}

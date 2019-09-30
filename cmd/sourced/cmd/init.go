package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
	"gopkg.in/src-d/go-cli.v0"

	"github.com/pkg/errors"
)

type initCmd struct {
	cli.PlainCommand `name:"init" short-description:"Initialize source{d} to work on local or GitHub orgs datasets" long-description:"Initialize source{d} to work on local or Github orgs datasets"`
}

type initLocalCmd struct {
	Command `name:"local" short-description:"Initialize source{d} to analyze local repositories" long-description:"Install, initialize, and start all the required docker containers, networks, volumes, and images.\n\nThe repos directory argument must point to a directory containing git repositories.\nIf it's not provided, the current working directory will be used."`

	Args struct {
		Reposdir string `positional-arg-name:"workdir"`
	} `positional-args:"yes"`
}

func (c *initLocalCmd) Execute(args []string) error {
	wdHandler, err := workdir.NewHandler()
	if err != nil {
		return err
	}

	reposdir, err := c.reposdirArg()
	if err != nil {
		return err
	}

	wd, err := workdir.InitLocal(reposdir)
	if err != nil {
		return err
	}

	if err := activate(wdHandler, wd); err != nil {
		return err
	}

	return OpenUI(60 * time.Minute)
}

func (c *initLocalCmd) reposdirArg() (string, error) {
	reposdir := c.Args.Reposdir
	reposdir = strings.TrimSpace(reposdir)

	var err error
	if reposdir == "" {
		reposdir, err = os.Getwd()
	} else {
		reposdir, err = filepath.Abs(reposdir)
	}

	if err != nil {
		return "", errors.Wrap(err, "could not get directory")
	}

	info, err := os.Stat(reposdir)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("path '%s' is not a valid directory", reposdir)
	}

	return reposdir, nil
}

type initOrgsCmd struct {
	Command `name:"orgs" short-description:"Initialize source{d} to analyze GitHub organizations" long-description:"Install, initialize, and start all the required docker containers, networks, volumes, and images.\n\nThe orgs argument must a comma-separated list of GitHub organization names to be analyzed."`

	Token     string `short:"t" long:"token" env:"SOURCED_GITHUB_TOKEN" description:"GitHub token for the passed organizations. It should be granted with 'repo' and 'read:org' scopes." required:"true"`
	WithForks bool   `long:"with-forks" description:"Download GitHub forked repositories"`
	Args      struct {
		Orgs []string `required:"yes"`
	} `positional-args:"yes" required:"1"`
}

func (c *initOrgsCmd) Execute(args []string) error {
	wdHandler, err := workdir.NewHandler()
	if err != nil {
		return err
	}

	orgs := c.orgsList()
	if err := c.validate(orgs); err != nil {
		return err
	}

	wd, err := workdir.InitOrgs(orgs, c.Token, c.WithForks)
	if err != nil {
		return err
	}

	if err := activate(wdHandler, wd); err != nil {
		return err
	}

	return OpenUI(60 * time.Minute)
}

// allows to pass organizations separated not only by a space
// but by comma as well
func (c *initOrgsCmd) orgsList() []string {
	orgs := c.Args.Orgs
	if len(c.Args.Orgs) == 1 {
		orgs = strings.Split(c.Args.Orgs[0], ",")
	}

	for i, org := range orgs {
		orgs[i] = strings.Trim(org, " ,")
	}

	return orgs
}

func (c *initOrgsCmd) validate(orgs []string) error {
	client := &http.Client{Transport: &authTransport{token: c.Token}}
	r, err := client.Get("https://api.github.com/user")
	if err != nil {
		return errors.Wrapf(err, "could not validate user token")
	}
	if r.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("github token is not valid")
	}

	for _, org := range orgs {
		r, err := client.Get("https://api.github.com/orgs/" + org)
		if err != nil {
			return errors.Wrapf(err, "could not validate organization")
		}
		if r.StatusCode == http.StatusNotFound {
			return fmt.Errorf("organization '%s' is not found", org)
		}
	}

	return nil
}

func activate(wdHandler *workdir.Handler, workdir *workdir.Workdir) error {
	// Before setting a new workdir, stop the current containers
	compose.Run(context.Background(), "stop")

	err := wdHandler.SetActive(workdir)
	if err != nil {
		return err
	}

	fmt.Printf("docker-compose working directory set to %s\n", workdir.Path)
	return compose.Run(context.Background(), "up", "--detach")
}

type authTransport struct {
	token string
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "token "+t.token)
	return http.DefaultTransport.RoundTrip(r)
}

func init() {
	c := rootCmd.AddCommand(&initCmd{})

	c.AddCommand(&initOrgsCmd{})
	c.AddCommand(&initLocalCmd{})
}

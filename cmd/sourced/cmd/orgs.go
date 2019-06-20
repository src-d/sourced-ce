package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/src-d/go-cli.v0"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
)

type orgsCmd struct {
	cli.PlainCommand `name:"orgs" short-description:"Manage services to analyze code from GitHub organizations" long-description:"Manage services to analyze code from GitHub organizations"`
}

type orgsInitCmd struct {
	Command `name:"init" short-description:"Install and initialize containers to analyze GitHub organizations" long-description:"Install, initialize, and start all the required docker containers, networks, volumes, and images.\n\nThe orgs argument must a comma-separated list of GitHub organization names to be analyzed."`

	Token string `short:"t" long:"token" env:"SOURCED_GITHUB_TOKEN" description:"Github token for the passed organizations. It should be granted with 'repo' and 'read:org' scopes." required:"true"`
	Args  struct {
		Orgs []string `required:"yes"`
	} `positional-args:"yes" required:"1"`
}

func (c *orgsInitCmd) Execute(args []string) error {
	orgs := c.orgsList()
	if err := c.validate(orgs); err != nil {
		return err
	}

	dir, err := workdir.InitWithOrgs(orgs, c.Token)
	if err != nil {
		return err
	}

	// Before setting a new workdir, stop the current containers
	compose.Run(context.Background(), "stop")

	err = workdir.SetActive(dir)
	if err != nil {
		return err
	}

	fmt.Printf("docker-compose working directory set to %s\n", strings.Join(orgs, ","))

	if err := compose.Run(context.Background(), "up", "--detach"); err != nil {
		return err
	}

	return OpenUI(60 * time.Minute)
}

// allows to pass organizations separated not only by a space
// but by comma as well
func (c *orgsInitCmd) orgsList() []string {
	orgs := c.Args.Orgs
	if len(c.Args.Orgs) == 1 {
		orgs = strings.Split(c.Args.Orgs[0], ",")
	}

	for i, org := range orgs {
		orgs[i] = strings.Trim(org, " ,")
	}

	return orgs
}

func (c *orgsInitCmd) validate(orgs []string) error {
	client := &http.Client{Transport: &authTransport{token: c.Token}}
	r, err := client.Get("https://api.github.com/user")
	if err != nil {
		return err
	}
	if r.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("github token is not valid")
	}

	for _, org := range orgs {
		r, err := client.Get("https://api.github.com/orgs/" + org)
		if err != nil {
			return err
		}
		if r.StatusCode == http.StatusNotFound {
			return fmt.Errorf("organization '%s' is not found", org)
		}
	}

	return nil
}

type authTransport struct {
	token string
}

func (t *authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "token "+t.token)
	return http.DefaultTransport.RoundTrip(r)
}

func init() {
	c := rootCmd.AddCommand(&orgsCmd{})
	c.AddCommand(&orgsInitCmd{})
}

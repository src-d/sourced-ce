package workdir

import (
	"encoding/base64"
	"sort"
	"strings"
)

// InitWithOrgs initialize workdir with remote list of organizations
func InitWithOrgs(orgs []string, token string) (string, error) {
	// be indifferent to the order of passed organizations
	sort.Strings(orgs)

	workdir := base64.URLEncoding.EncodeToString([]byte(strings.Join(orgs, ",")))
	workdirPath, err := absolutePath(workdir)
	if err != nil {
		return "", err
	}

	envf := envFile{
		Workdir:             workdir,
		GithubOrganizations: orgs,
		GithubToken:         token,
	}
	if err := initWorkdir(workdirPath, envf); err != nil {
		return "", err
	}

	return workdir, nil
}

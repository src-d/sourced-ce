package workdir

import (
	"encoding/base64"
	"strings"
)

// InitWithOrgs initialize workdir with remote list of organizations
func InitWithOrgs(orgs []string, token string) (string, error) {
	workdir := base64.StdEncoding.EncodeToString([]byte(strings.Join(orgs, ",")))
	workdirPath, err := absolutePath(workdir)
	if err != nil {
		return "", err
	}

	envf := envFile{
		Workdir:             workdir,
		GithubOrganizations: orgs,
		GithubToken:         token,
	}
	if err := initWorkdir(workdirPath, envf.String); err != nil {
		return "", err
	}

	return workdir, nil
}

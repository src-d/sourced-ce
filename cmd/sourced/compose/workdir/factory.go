package workdir

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	composefile "github.com/src-d/sourced-ce/cmd/sourced/compose/file"
)

// InitLocal initializes the workdir for local path and returns the Workdir instance
func InitLocal(reposdir string) (*Workdir, error) {
	dirName := encodeDirName(reposdir)

	envf := envFile{
		Workdir:  dirName,
		ReposDir: reposdir,
	}

	return initialize(dirName, "local", envf)
}

// InitOrgs initializes the workdir for organizations and returns the Workdir instance
func InitOrgs(orgs []string, token string) (*Workdir, error) {
	// be indifferent to the order of passed organizations
	sort.Strings(orgs)
	dirName := encodeDirName(strings.Join(orgs, ","))

	envf := envFile{
		Workdir:             dirName,
		GithubOrganizations: orgs,
		GithubToken:         token,
	}

	return initialize(dirName, "orgs", envf)
}

func encodeDirName(dirName string) string {
	return base64.URLEncoding.EncodeToString([]byte(dirName))
}

func buildAbsPath(dirName, subPath string) (string, error) {
	path, err := workdirsPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, subPath, dirName), nil
}

func initialize(dirName string, subPath string, envf envFile) (*Workdir, error) {
	path, err := workdirsPath()
	if err != nil {
		return nil, err
	}

	workdir := filepath.Join(path, subPath, dirName)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "could not create working directory")
	}

	defaultFilePath, err := composefile.InitDefault()
	if err != nil {
		return nil, err
	}

	composePath := filepath.Join(workdir, "docker-compose.yml")
	if err := link(defaultFilePath, composePath); err != nil {
		return nil, err
	}

	defaultOverridePath, err := composefile.InitDefaultOverride()
	if err != nil {
		return nil, err
	}

	workdirOverridePath := filepath.Join(workdir, "docker-compose.override.yml")
	if err := link(defaultOverridePath, workdirOverridePath); err != nil {
		return nil, err
	}

	envPath := filepath.Join(workdir, ".env")
	contents := envf.String()
	err = ioutil.WriteFile(envPath, []byte(contents), 0644)

	if err != nil {
		return nil, errors.Wrap(err, "could not write .env file")
	}

	b := &builder{workdirsPath: path}
	return b.build(workdir)
}

type envFile struct {
	Workdir             string
	ReposDir            string
	GithubOrganizations []string
	GithubToken         string
}

func (f *envFile) String() string {
	volumeType := "bind"
	volumeSource := f.ReposDir
	gitbaseSiva := ""
	if f.ReposDir == "" {
		volumeType = "volume"
		volumeSource = "gitbase_repositories"
		gitbaseSiva = "true"
	}

	return fmt.Sprintf(`COMPOSE_PROJECT_NAME=srcd-%s
	GITBASE_VOLUME_TYPE=%s
	GITBASE_VOLUME_SOURCE=%s
	GITBASE_SIVA=%s
	GITHUB_ORGANIZATIONS=%s
	GITHUB_TOKEN=%s
	`, f.Workdir, volumeType, volumeSource, gitbaseSiva,
		strings.Join(f.GithubOrganizations, ","), f.GithubToken)
}

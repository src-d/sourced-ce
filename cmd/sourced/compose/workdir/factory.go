package workdir

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/pbnjay/memory"
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
func InitOrgs(orgs []string, token string, withForks bool) (*Workdir, error) {
	// be indifferent to the order of passed organizations
	sort.Strings(orgs)
	dirName := encodeDirName(strings.Join(orgs, ","))

	envf := envFile{}
	err := readEnvFile(dirName, "orgs", &envf)
	if err == nil && envf.WithForks != withForks {
		return nil, ErrInitFailed.Wrap(
			fmt.Errorf("workdir was previously initialized with a different value for forks support"))
	}
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// re-create env file to make sure all fields are updated
	envf = envFile{
		Workdir:             dirName,
		GithubOrganizations: orgs,
		GithubToken:         token,
		WithForks:           withForks,
	}

	return initialize(dirName, "orgs", envf)
}

func readEnvFile(dirName string, subPath string, envf *envFile) error {
	workdir, err := workdirPath(dirName, subPath)
	if err != nil {
		return err
	}

	envPath := filepath.Join(workdir, ".env")
	b, err := ioutil.ReadFile(envPath)
	if err != nil {
		return err
	}

	return envf.UnmarshalEnv(b)
}

func encodeDirName(dirName string) string {
	return base64.URLEncoding.EncodeToString([]byte(dirName))
}

func workdirPath(dirName string, subPath string) (string, error) {
	path, err := workdirsPath()
	if err != nil {
		return "", err
	}

	workdir := filepath.Join(path, subPath, dirName)
	if err != nil {
		return "", err
	}

	return workdir, nil
}

func initialize(dirName string, subPath string, envf envFile) (*Workdir, error) {
	path, err := workdirsPath()
	if err != nil {
		return nil, err
	}

	workdir, err := workdirPath(dirName, subPath)
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
	contents, err := envf.MarshalEnv()
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(envPath, contents, 0644)

	if err != nil {
		return nil, errors.Wrap(err, "could not write .env file")
	}

	b := &builder{workdirsPath: path}
	return b.Build(workdir)
}

type envFile struct {
	Workdir             string
	ReposDir            string
	GithubOrganizations []string
	GithubToken         string
	WithForks           bool
}

func (f *envFile) MarshalEnv() ([]byte, error) {
	volumeType := "bind"
	volumeSource := f.ReposDir
	gitbaseSiva := ""
	if f.ReposDir == "" {
		volumeType = "volume"
		volumeSource = "gitbase_repositories"
		gitbaseSiva = "true"
	}

	noForks := "true"
	if f.WithForks {
		noForks = "false"
	}

	// limit CPU for containers
	gitbaseLimitCPU := "0.0"
	gitcollectorLimitCPU := "0.0"
	dockerCPUs, err := dockerNumCPU()
	if err != nil { // show warning
		fmt.Println(err)
	}
	// apply gitbase resource limits only when docker runs without any global limits
	// it's default behaviour on linux
	if runtime.NumCPU() == dockerCPUs {
		gitbaseLimitCPU = fmt.Sprintf("%.1f", float32(dockerCPUs)-0.1)
	}
	if dockerCPUs > 0 {
		halfCPUs := float32(dockerCPUs)/2.0
		if halfCPUs < 1 {
			halfCPUs = 1
		}
		gitcollectorLimitCPU = fmt.Sprintf("%.1f", halfCPUs-0.1)
	}

	// limit memory for containers
	gitbaseLimitMem := "0"
	dockerMem, err := dockerTotalMem()
	if err != nil { // show warning
		fmt.Println(err)
	}
	// apply memory limits only when only when docker runs without any global limits
	// it's default behaviour on linux
	if dockerMem == memory.TotalMemory() {
		gitbaseLimitMem = strconv.FormatUint(uint64(float64(dockerMem)*0.9), 10)
	}

	result := fmt.Sprintf(`COMPOSE_PROJECT_NAME=srcd-%s
GITBASE_VOLUME_TYPE=%s
GITBASE_VOLUME_SOURCE=%s
GITBASE_SIVA=%s
GITHUB_ORGANIZATIONS=%s
GITHUB_TOKEN=%s
NO_FORKS=%s
GITBASE_LIMIT_CPU=%s
GITCOLLECTOR_LIMIT_CPU=%s
GITBASE_LIMIT_MEM=%s
`, f.Workdir, volumeType, volumeSource, gitbaseSiva,
		strings.Join(f.GithubOrganizations, ","), f.GithubToken, noForks,
		gitbaseLimitCPU, gitcollectorLimitCPU, gitbaseLimitMem)

	return []byte(result), nil
}

func (f *envFile) UnmarshalEnv(b []byte) error {
	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		value := parts[1]

		switch parts[0] {
		case "COMPOSE_PROJECT_NAME":
			f.Workdir = strings.TrimPrefix(value, "srcd-")
		case "GITBASE_VOLUME_SOURCE":
			f.ReposDir = value
		case "GITHUB_ORGANIZATIONS":
			if value != "" {
				f.GithubOrganizations = strings.Split(value, ",")
			}
		case "GITHUB_TOKEN":
			f.GithubToken = value
		case "NO_FORKS":
			if value == "false" {
				f.WithForks = true
			}
		}
	}

	if len(f.GithubOrganizations) > 0 {
		f.ReposDir = ""
	}

	return scanner.Err()
}

// returns number of CPUs available to docker
func dockerNumCPU() (int, error) {
	// use cli instead of connection to docker server directly
	// in case server exposed by http or non-default socket path
	info, err := exec.Command("docker", "info", "--format", "{{.NCPU}}").Output()
	if err != nil {
		return 0, err
	}

	cpus, err := strconv.Atoi(strings.TrimSpace(string(info)))
	if err != nil || cpus == 0 {
		return 0, fmt.Errorf("Couldn't get number of available CPUs in docker")
	}

	return cpus, nil
}

// returns total memory in bytes available to docker
func dockerTotalMem() (uint64, error) {
	info, err := exec.Command("docker", "info", "--format", "{{.MemTotal}}").Output()
	if err != nil {
		return 0, err
	}

	mem, err := strconv.ParseUint(strings.TrimSpace(string(info)), 10, 64)
	if err != nil || mem == 0 {
		return 0, fmt.Errorf("Couldn't get of available memory in docker")
	}

	return mem, nil
}

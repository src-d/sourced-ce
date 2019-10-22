package compose

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/blang/semver"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
	"github.com/src-d/sourced-ce/cmd/sourced/dir"

	"github.com/pkg/errors"
	goerrors "gopkg.in/src-d/go-errors.v1"
)

// v1.20.0 is the first version that supports `--compatibility` flag we rely on
// there is no mention of it in changelog
// the version has been found by trying downgrading unless it started to error
var minDockerComposeVersion = semver.Version{
	Major: 1,
	Minor: 20,
	Patch: 0,
}

// this version is choosen to be always compatible with docker-compose version
// docker-compose v1.20.0 introduced compose files version 3.6
// which requires Docker Engine 18.02.0 or above
var minDockerVersion = semver.Version{
	Major: 18,
	Minor: 2,
	Patch: 0,
}

// dockerComposeVersion is the version of docker-compose to download
// if docker-compose isn't already present in the system
const dockerComposeVersion = "1.24.0"

var composeContainerURL = fmt.Sprintf("https://github.com/docker/compose/releases/download/%s/run.sh", dockerComposeVersion)

// ErrComposeAlternative is returned when docker-compose alternative could not be installed
var ErrComposeAlternative = goerrors.NewKind("error while trying docker-compose container alternative")

type Compose struct {
	bin            string
	workdirHandler *workdir.Handler
}

func (c *Compose) Run(ctx context.Context, arg ...string) error {
	return c.RunWithIO(ctx, os.Stdin, os.Stdout, os.Stderr, arg...)
}

func (c *Compose) RunWithIO(ctx context.Context, stdin io.Reader,
	stdout, stderr io.Writer, arg ...string) error {
	arg = append([]string{"--compatibility"}, arg...)
	cmd := exec.CommandContext(ctx, c.bin, arg...)

	wd, err := c.workdirHandler.Active()
	if err != nil {
		return err
	}

	if err := c.workdirHandler.Validate(wd); err != nil {
		return err
	}

	cmd.Dir = wd.Path
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func newCompose() (*Compose, error) {
	// check docker first and exit fast
	dockerVersion, err := getDockerVersion()
	if err != nil {
		return nil, err
	}
	if !dockerVersion.GE(minDockerVersion) {
		return nil, fmt.Errorf("Minimal required docker version is %s but %s found", minDockerVersion, dockerVersion)
	}

	workdirHandler, err := workdir.NewHandler()
	if err != nil {
		return nil, err
	}

	bin, err := getOrInstallComposeBinary()
	if err != nil {
		return nil, err
	}
	dockerComposeVersion, err := getDockerComposeVersion(bin)
	if err != nil {
		return nil, err
	}
	if !dockerComposeVersion.GE(minDockerComposeVersion) {
		return nil, fmt.Errorf("Minimal required docker-compose version is %s but %s found", minDockerComposeVersion, dockerComposeVersion)
	}

	return &Compose{
		bin:            bin,
		workdirHandler: workdirHandler,
	}, nil
}

func getOrInstallComposeBinary() (string, error) {
	path, err := exec.LookPath("docker-compose")
	if err == nil {
		bin := strings.TrimSpace(path)
		if bin != "" {
			return bin, nil
		}
	}

	path, err = getOrInstallComposeContainer()
	if err != nil {
		return "", ErrComposeAlternative.Wrap(err)
	}

	return path, nil
}

func getOrInstallComposeContainer() (altPath string, err error) {
	datadir, err := dir.Path()
	if err != nil {
		return "", err
	}

	dirPath := filepath.Join(datadir, "bin")
	path := filepath.Join(dirPath, fmt.Sprintf("docker-compose-%s.sh", dockerComposeVersion))

	readExecAccessMode := os.FileMode(0500)

	if info, err := os.Stat(path); err == nil {
		if info.Mode()&readExecAccessMode != readExecAccessMode {
			return "", fmt.Errorf("%s can not be run", path)
		}

		return path, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	if err := downloadCompose(path); err != nil {
		return "", err
	}

	cmd := exec.CommandContext(context.Background(), "chmod", "+x", path)
	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "cannot change permission to %s", path)
	}

	return path, nil
}

func downloadCompose(path string) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("compose in container is not compatible with Windows")
	}

	return dir.DownloadURL(composeContainerURL, path)
}

func Run(ctx context.Context, arg ...string) error {
	comp, err := newCompose()
	if err != nil {
		return err
	}

	return comp.Run(ctx, arg...)
}

func RunWithIO(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, arg ...string) error {
	comp, err := newCompose()
	if err != nil {
		return err
	}

	return comp.RunWithIO(ctx, stdin, stdout, stderr, arg...)
}

var dockerVersionRe = regexp.MustCompile(`version (\d+).(\d+).(\d+)`)
var dockerComposeVersionRe = regexp.MustCompile(`version (\d+.\d+.\d+)`)

// docker doesn't use semver, so simple `semver.Parse` would fail
// but semver.Version struct fits us to allow simple comparation
func getDockerVersion() (*semver.Version, error) {
	if _, err := exec.LookPath("docker"); err != nil {
		return nil, err
	}

	out, err := exec.Command("docker", "--version").Output()
	if err != nil {
		return nil, err
	}

	submatches := dockerVersionRe.FindSubmatch(out)
	if len(submatches) != 4 {
		return nil, fmt.Errorf("can't parse docker version")
	}

	v := &semver.Version{}
	v.Major, err = strconv.ParseUint(string(submatches[1]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse docker version")
	}
	v.Minor, err = strconv.ParseUint(string(submatches[2]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse docker version")
	}
	v.Patch, err = strconv.ParseUint(string(submatches[3]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse docker version")
	}

	return v, nil
}

func getDockerComposeVersion(bin string) (*semver.Version, error) {
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		return nil, err
	}

	submatches := dockerComposeVersionRe.FindSubmatch(out)
	if len(submatches) != 2 {
		return nil, fmt.Errorf("can't parse docker-compose version")
	}

	v, err := semver.ParseTolerant(string(submatches[1]))
	if err != nil {
		return nil, fmt.Errorf("can't parse docker-compose version: %s", err)
	}

	return &v, nil
}

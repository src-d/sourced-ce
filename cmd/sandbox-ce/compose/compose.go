package compose

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	composefile "github.com/src-d/superset-compose/cmd/sandbox-ce/compose/file"
	"github.com/src-d/superset-compose/cmd/sandbox-ce/dir"
	datadir "github.com/src-d/superset-compose/cmd/sandbox-ce/dir"

	"github.com/pkg/errors"
)

// composeContainerURL is the url docker-compose is downloaded from in case
// that it's not already present in the system
const composeContainerURL = "https://github.com/docker/compose/releases/download/1.24.0/run.sh"

var envKeys = []string{"GITBASE_REPOS_DIR"}

const activeWorkdir = "__active__"

type Compose struct {
	bin string
}

func (c *Compose) Run(ctx context.Context, arg ...string) error {
	return c.RunWithIO(ctx, os.Stdin, os.Stdout, os.Stderr, arg...)
}

func (c *Compose) RunWithIO(ctx context.Context, stdin io.Reader,
	stdout, stderr io.Writer, arg ...string) error {
	cmd := exec.CommandContext(ctx, c.bin, arg...)

	dir, err := workdirPath(activeWorkdir)
	if err != nil {
		return err
	}

	cmd.Dir = dir

	var compOpts []string
	for _, key := range envKeys {
		val, ok := os.LookupEnv(key)
		if !ok {
			continue
		}

		compOpts = append(compOpts, fmt.Sprintf("-e %s=%s", key, val))
	}

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("COMPOSE_OPTIONS=%s", strings.Join(compOpts, " ")))
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func NewCompose() (*Compose, error) {
	bin, err := getOrInstallComposeBinary()
	if err != nil {
		return nil, err
	}

	return &Compose{bin: bin}, nil
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
		return "", errors.Wrapf(err, "error while getting docker-compose container")
	}

	return path, nil
}

func getOrInstallComposeContainer() (string, error) {
	datadir, err := dir.Path()
	if err != nil {
		return "", err
	}

	dirPath := filepath.Join(datadir, "bin")
	path := filepath.Join(dirPath, "docker-compose")

	if _, err := os.Stat(path); err == nil {
		return path, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", errors.Wrapf(err, "error while creating directory %s", dirPath)
	}

	if err := downloadCompose(path); err != nil {
		return "", errors.Wrapf(err, "error downloading %s", composeContainerURL)
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
	comp, err := NewCompose()
	if err != nil {
		return err
	}

	return comp.Run(ctx, arg...)
}

func RunWithIO(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, arg ...string) error {
	comp, err := NewCompose()
	if err != nil {
		return err
	}

	return comp.RunWithIO(ctx, stdin, stdout, stderr, arg...)
}

// InitWorkdir creates a working directory in ~/.srcd for the given repositories
// directory. The working directory will contain a docker-compose.yml and a
// .env file.
// If the directory is already initialized the function returns with no error
func InitWorkdir(reposdir string) (string, error) {
	defaultFilePath, err := composefile.InitDefault()
	if err != nil {
		return "", err
	}

	workdir, err := workdirPath(reposdir)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return "", errors.Wrap(err, "could not create working directory")
	}

	composePath := filepath.Join(workdir, "docker-compose.yml")
	_, err = os.Stat(composePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", errors.Wrap(err, "could not read the existing docker-compose.yml file")
		}

		err = os.Symlink(defaultFilePath, composePath)
		if err != nil {
			return "", errors.Wrap(err, "could not create symlink to docker-compose.yml file")
		}
	}

	envPath := filepath.Join(workdir, ".env")
	emptyFile, err := isEmptyFile(envPath)
	if err != nil {
		return "", errors.Wrap(err, "could not read .env file contents")
	}

	if emptyFile {
		hash := sha1.Sum([]byte(reposdir))
		hashSt := hex.EncodeToString(hash[:])

		contents := fmt.Sprintf(
			`GITBASE_REPOS_DIR=%s
COMPOSE_PROJECT_NAME=srcd-%s
`,
			reposdir, hashSt)

		err = ioutil.WriteFile(envPath, []byte(contents), 0644)
		if err != nil {
			return "", errors.Wrap(err, "could not write .env file")
		}
	}

	return workdir, nil
}

// isEmptyFile returns true if the file does not exist or if it exists but
// contains empty text
func isEmptyFile(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}

		return true, nil
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	strContents := string(contents)
	return strings.TrimSpace(strContents) == "", nil
}

// SetActiveWorkdir creates a symlink from the fixed active workdir path
// to the given workdir. The workdir should be the path returned by InitWorkdir
func SetActiveWorkdir(workdir string) error {
	dir, err := workdirPath(activeWorkdir)
	if err != nil {
		return err
	}

	_, err = os.Stat(dir)
	if !os.IsNotExist(err) {
		err = os.Remove(dir)
		if err != nil {
			return errors.Wrap(err, "could not delete the previous active workdir directory symlink")
		}
	}

	return os.Symlink(workdir, dir)
}

// workdirPath returns the absolute path to
// $HOME/.srcd/workdirs/reposdir
func workdirPath(reposdir string) (string, error) {
	path, err := datadir.Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "workdirs", reposdir), nil
}

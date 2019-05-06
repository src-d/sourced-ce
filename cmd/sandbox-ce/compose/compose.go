package compose

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// composeContainerPath is the url docker-compose is downloaded from in case
// that it's not already present in the system
const composeContainerPath = "https://github.com/docker/compose/releases/download/1.24.0/run.sh"

var envKeys = map[string]bool{
	"GITBASE_REPOS_DIR": true,
}

type Compose struct {
	bin string
}

func (c *Compose) Run(ctx context.Context, arg ...string) error {
	cmd := exec.CommandContext(ctx, c.bin, arg...)

	var compOpts []string
	for key := range envKeys {
		val, ok := os.LookupEnv(key)
		if !ok {
			continue
		}

		compOpts = append(compOpts, fmt.Sprintf("-e %s=%s", key, val))
	}

	cmd.Env = []string{
		fmt.Sprintf("COMPOSE_OPTIONS=%s", strings.Join(compOpts, " ")),
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
	homedir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrapf(err, "unable to get home dir")
	}

	dirPath := filepath.Join(homedir, ".srcd", "bin")
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
		return "", errors.Wrapf(err, "error downloading %s", composeContainerPath)
	}

	cmd := exec.CommandContext(context.Background(), "chmod", "+x", path)
	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "cannot change permission to %s", path)
	}

	return path, nil
}

func downloadCompose(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(composeContainerPath)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Run(ctx context.Context, arg ...string) error {
	comp, err := NewCompose()
	if err != nil {
		return err
	}

	return comp.Run(ctx, arg...)
}

// Package workdir provides functions to manage docker compose working
// directories inside the $HOME/.srcd/workdirs directory
package workdir

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	composefile "github.com/src-d/sourced-ce/cmd/sandbox-ce/compose/file"
	datadir "github.com/src-d/sourced-ce/cmd/sandbox-ce/dir"

	"github.com/pkg/errors"
)

const activeDir = "__active__"

// Init creates a working directory in ~/.srcd for the given repositories
// directory. The working directory will contain a docker-compose.yml and a
// .env file.
// If the directory is already initialized the function returns with no error.
// The returned value is the absolute path to $HOME/.srcd/workdirs/reposdir
func Init(reposdir string) (string, error) {
	defaultFilePath, err := composefile.InitDefault()
	if err != nil {
		return "", err
	}

	workdir, err := path(reposdir)
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

// SetActive creates a symlink from the fixed active workdir path
// to the workdir for the given repos dir.
func SetActive(reposdir string) error {
	dir, err := path(activeDir)
	if err != nil {
		return err
	}

	workdir, err := path(reposdir)
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

// Active returns th absolute path to $HOME/.srcd/workdirs/__active__
func Active() (string, error) {
	return path(activeDir)
}

// path returns the absolute path to
// $HOME/.srcd/workdirs/reposdir
func path(reposdir string) (string, error) {
	path, err := datadir.Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "workdirs", reposdir), nil
}

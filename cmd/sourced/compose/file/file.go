// Package file provides functions to manage docker compose files inside the
// $HOME/.sourced/compose-files directory
package file

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	datadir "github.com/src-d/sourced-ce/cmd/sourced/dir"

	"github.com/pkg/errors"
	goerrors "gopkg.in/src-d/go-errors.v1"
)

// ErrConfigDownload is returned when docker-compose.yml could not be downloaded
var ErrConfigDownload = goerrors.NewKind("docker-compose.yml config file could not be downloaded")

// ErrConfigActivation is returned when docker-compose.yml could not be set as active
var ErrConfigActivation = goerrors.NewKind("docker-compose.yml could not be set as active")

const (
	orgName         = "src-d"
	repoName        = "sourced-ce"
	composeFileTmpl = "https://raw.githubusercontent.com/%s/%s/%s/docker-compose.yml"
)

var version = "master"

// activeDir is the name of the directory containing the symlink to the
// active docker compose file
const activeDir = "__active__"

// RevOrURL is a revision (tag name, full sha1) or a valid URL to a
// docker-compose.yml file
type RevOrURL = string

// composeFileURL returns the URL to download the raw github docker-compose.yml
// file for the given revision (tag or full sha1)
func composeFileURL(revision string) string {
	return fmt.Sprintf(composeFileTmpl, orgName, repoName, revision)
}

// SetVersion sets the version rewritten by the CI build
func SetVersion(v string) {
	version = v
}

// InitDefault checks if there is an active docker compose file, and if there
// isn't the file for this release is downloaded.
// The current build version must be set with SetVersion.
// It returns the absolute path to the active docker-compose.yml file
func InitDefault() (string, error) {
	activeFilePath, err := path(activeDir)
	if err != nil {
		return "", err
	}

	_, err = os.Stat(activeFilePath)
	if err == nil {
		return activeFilePath, nil
	}

	if !os.IsNotExist(err) {
		return "", err
	}

	err = ActivateFromRemote(version)
	if err != nil {
		return "", err
	}

	return activeFilePath, nil
}

// ActivateFromRemote downloads the docker-compose.yml file from the given revision
// or URL, and sets it as the active compose file.
func ActivateFromRemote(revOrURL RevOrURL) (err error) {
	var url string
	if isURL(revOrURL) {
		url = revOrURL
	} else {
		url = composeFileURL(revOrURL)
	}

	outPath, err := path(revOrURL)
	if err != nil {
		return err
	}

	err = datadir.DownloadURL(url, outPath)
	if err != nil {
		return ErrConfigDownload.Wrap(err)
	}

	err = SetActive(revOrURL)
	if err != nil {
		return ErrConfigActivation.Wrap(err)
	}

	return nil
}

// SetActive makes a symlink from
// $HOME/.sourced/compose-files/__active__/docker-compose.yml to the compose file
// for the given revision or URL.
func SetActive(revOrURL RevOrURL) error {
	filePath, err := path(revOrURL)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return errors.Wrapf(err, "could not find a docker-compose.yml file in `%s`", filePath)
		}

		return err
	}

	activeFilePath, err := path(activeDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(activeFilePath), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "error while creating directory for %s", activeFilePath)
	}

	if _, err := os.Lstat(activeFilePath); err == nil {
		if err := os.Remove(activeFilePath); err != nil {
			return errors.Wrap(err, "failed to unlink")
		}
	}

	return os.Symlink(filePath, activeFilePath)
}

// Active returns the revision (tag name, full sha1) or the URL of the active
// docker compose file
func Active() (RevOrURL, error) {
	activeFilePath, err := path(activeDir)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(activeFilePath); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", err
	}

	dest, err := filepath.EvalSymlinks(activeFilePath)
	if err != nil {
		return "", err
	}

	_, name := filepath.Split(filepath.Dir(dest))
	return composeName(name), nil
}

// List returns a list of installed docker compose files. Each name is the
// revision (tag name, full sha1) or the URL
func List() ([]RevOrURL, error) {
	list := []RevOrURL{}

	dir, err := dir()
	if err != nil {
		return list, err
	}

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return list, nil
		}

		return list, err
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return list, err
	}

	for _, f := range files {
		entry := f.Name()
		if entry == activeDir {
			continue
		}

		list = append(list, composeName(entry))
	}

	return list, nil
}

func composeName(rev string) string {
	if decoded, err := base64.URLEncoding.DecodeString(rev); err == nil {
		return string(decoded)
	}

	return rev
}

func isURL(revOrURL RevOrURL) bool {
	_, err := url.ParseRequestURI(revOrURL)
	return err == nil
}

// dir returns the absolute path for $HOME/.sourced/compose-files
func dir() (string, error) {
	path, err := datadir.Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "compose-files"), nil
}

// path returns the absolute path to
// $HOME/.sourced/compose-files/revOrURL/docker-compose.yml
func path(revOrURL RevOrURL) (string, error) {
	composeDirPath, err := dir()
	if err != nil {
		return "", err
	}

	subPath := revOrURL
	if isURL(revOrURL) {
		subPath = base64.URLEncoding.EncodeToString([]byte(revOrURL))
	}

	dirPath := filepath.Join(composeDirPath, subPath)

	return filepath.Join(dirPath, "docker-compose.yml"), nil
}

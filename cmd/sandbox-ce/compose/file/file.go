// Package file provides functions to manage docker compose files inside the
// $HOME/.srcd/compose-files directory
package file

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	datadir "github.com/src-d/superset-compose/cmd/sandbox-ce/dir"

	"github.com/pkg/errors"
)

const (
	orgName         = "src-d"
	repoName        = "superset-compose"
	composeFileTmpl = "https://raw.githubusercontent.com/%s/%s/%s/docker-compose.yml"
)

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

// InitDefault downloads the master docker-compose.yml only if it does not
// exist. It returns the absolute path to the active docker-compose.yml file
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

	err = Download("master")
	if err != nil {
		return "", nil
	}

	return activeFilePath, nil
}

// Download downloads the docker-compose.yml file from the given revision
// or URL. The file is set as the active compose file.
func Download(revOrURL RevOrURL) error {
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
		return err
	}

	return SetActive(revOrURL)
}

// SetActive makes a symlink from
// $HOME/.srcd/compose-files/__active__/docker-compose.yml to the compose file
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

	_, rev := filepath.Split(filepath.Dir(dest))
	return rev, nil
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

		name := entry

		decoded, err := base64.StdEncoding.DecodeString(entry)
		if err == nil {
			name = string(decoded)
		}

		list = append(list, name)
	}

	return list, nil
}

func isURL(revOrURL RevOrURL) bool {
	_, err := url.ParseRequestURI(revOrURL)
	return err == nil
}

// dir returns the absolute path for $HOME/.srcd/compose-files
func dir() (string, error) {
	path, err := datadir.Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "compose-files"), nil
}

// path returns the absolute path to
// $HOME/.srcd/compose-files/revOrURL/docker-compose.yml
func path(revOrURL RevOrURL) (string, error) {
	composeDirPath, err := dir()
	if err != nil {
		return "", err
	}

	subPath := revOrURL
	if isURL(revOrURL) {
		subPath = base64.StdEncoding.EncodeToString([]byte(revOrURL))
	}

	dirPath := filepath.Join(composeDirPath, subPath)

	return filepath.Join(dirPath, "docker-compose.yml"), nil
}

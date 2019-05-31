// Package dir provides functions to manage the $HOME/.srcd directory
package dir

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Path returns the absolute path for $HOME/.srcd
func Path() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "could not detect home directory")
	}

	return filepath.Join(homedir, ".srcd"), nil
}

// DownloadURL downloads the given url to a file to the
// dst path, creating the directory if it's needed
func DownloadURL(url, dst string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %v", resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

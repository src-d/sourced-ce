// Package dir provides functions to manage the $HOME/.sourced and /tmp/srcd
// directories
package dir

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	goerrors "gopkg.in/src-d/go-errors.v1"
)

// ErrNotExist is returned when .sourced dir does not exists
var ErrNotExist = goerrors.NewKind("%s does not exist")

// Path returns the absolute path for $SOURCED_DIR, or $HOME/.sourced if unset
func Path() (string, error) {
	if d := os.Getenv("SOURCED_DIR"); d != "" {
		return filepath.Abs(d)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "could not detect home directory")
	}

	srcdDir := filepath.Join(homedir, ".sourced")
	_, err = os.Lstat(srcdDir)
	if os.IsNotExist(err) {
		return "", ErrNotExist.New(srcdDir)
	}

	return srcdDir, nil
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

// TmpPath returns the absolute path for /tmp/srcd
func TmpPath() string {
	return filepath.Join(os.TempDir(), "srcd")
}

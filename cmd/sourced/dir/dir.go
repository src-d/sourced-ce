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

// ErrNotValid is returned when config dir is not valid
var ErrNotValid = goerrors.NewKind("%s is not a valid config directory: %s")

// Path returns the absolute path for $SOURCED_DIR, or $HOME/.sourced if unset
func Path() (string, error) {
	srcdDir, err := path()
	if err != nil {
		return "", err
	}

	if err := validate(srcdDir); err != nil {
		return "", err
	}

	return srcdDir, nil
}

func path() (string, error) {
	if d := os.Getenv("SOURCED_DIR"); d != "" {
		abs, err := filepath.Abs(d)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("could not resolve SOURCED_DIR='%s'", d))
		}

		return abs, nil
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "could not detect home directory")
	}

	return filepath.Join(homedir, ".sourced"), nil
}

// Prepare tries to create the config directory, returning an error if it could not
// be created, or nill if already exist or was successfully created.
func Prepare() error {
	srcdDir, err := path()
	if err != nil {
		return err
	}

	err = validate(srcdDir)
	if ErrNotExist.Is(err) {
		return os.MkdirAll(srcdDir, os.ModePerm)
	}

	return err
}

// validate validates that the passed config dir path is valid
func validate(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return ErrNotExist.New(path)
	}

	if err != nil {
		return ErrNotValid.New(path, err)
	}

	readWriteAccessMode := os.FileMode(0700)
	if info.Mode()&readWriteAccessMode != readWriteAccessMode {
		return ErrNotValid.New(path, "it has no read-write access")
	}

	if !info.IsDir() {
		return ErrNotValid.New(path, "it is not a directory")
	}

	return nil
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

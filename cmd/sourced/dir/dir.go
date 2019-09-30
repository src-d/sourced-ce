// Package dir provides functions to manage the config directories.
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

// ErrNotExist is returned when config dir does not exists
var ErrNotExist = goerrors.NewKind("%s does not exist")

// ErrNotValid is returned when config dir is not valid
var ErrNotValid = goerrors.NewKind("%s is not a valid config directory: %s")

// ErrNetwork is returned when could not download
var ErrNetwork = goerrors.NewKind("network error downloading %s")

// ErrWrite is returned when could not write
var ErrWrite = goerrors.NewKind("write error at %s")

// Path returns the absolute path for $SOURCED_DIR, or $HOME/.sourced if unset
// and returns an error if it does not exist or it could not be read.
func Path() (string, error) {
	srcdDir, err := srcdPath()
	if err != nil {
		return "", err
	}

	if err := validate(srcdDir); err != nil {
		return "", err
	}

	return srcdDir, nil
}

func srcdPath() (string, error) {
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
// be created, or nil if already exist or was successfully created.
func Prepare() error {
	srcdDir, err := srcdPath()
	if err != nil {
		return err
	}

	err = validate(srcdDir)
	if ErrNotExist.Is(err) {
		if err := os.MkdirAll(srcdDir, os.ModePerm); err != nil {
			return ErrNotValid.New(srcdDir, err)
		}

		return nil
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

	if !info.IsDir() {
		return ErrNotValid.New(path, "it is not a directory")
	}

	readWriteAccessMode := os.FileMode(0700)
	if info.Mode()&readWriteAccessMode != readWriteAccessMode {
		return ErrNotValid.New(path, "it has no read-write access")
	}

	return nil
}

// DownloadURL downloads the given url to a file to the
// dst path, creating the directory if it's needed
func DownloadURL(url, dst string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return ErrNetwork.Wrap(err, url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrNetwork.Wrap(fmt.Errorf("HTTP status %v", resp.Status), url)
	}

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return ErrWrite.Wrap(err, filepath.Dir(dst))
	}

	out, err := os.Create(dst)
	if err != nil {
		return ErrWrite.Wrap(err, dst)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return ErrWrite.Wrap(err, dst)
	}

	return nil
}

// TmpPath returns the absolute path for /tmp/srcd
func TmpPath() string {
	return filepath.Join(os.TempDir(), "srcd")
}

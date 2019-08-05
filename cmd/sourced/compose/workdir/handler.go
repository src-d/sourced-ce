package workdir

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Handler provides a way to interact with all the workdirs by exposing the following operations:
//   - read/set/unset active workdir,
//   - remove/validate a workdir,
//   - list workdirs.
type Handler struct {
	workdirsPath string
	builder      *builder
}

// NewHandler creates a handler that manages workdirs in the path returned by
// the `workdirsPath` function
func NewHandler() (*Handler, error) {
	path, err := workdirsPath()
	if err != nil {
		return nil, err
	}

	return &Handler{
		workdirsPath: path,
		builder:      &builder{workdirsPath: path},
	}, nil
}

// SetActive creates a symlink from the fixed active workdir path to the prodived workdir
func (h *Handler) SetActive(w *Workdir) error {
	path, err := h.activeAbsolutePath()
	if err != nil {
		return err
	}

	if err := h.UnsetActive(); err != nil {
		return err
	}

	err = os.Symlink(w.Path, path)
	if os.IsExist(err) {
		return nil
	}

	return err
}

// UnsetActive removes symlink for active workdir
func (h *Handler) UnsetActive() error {
	path, err := h.activeAbsolutePath()
	if err != nil {
		return err
	}

	_, err = os.Lstat(path)
	if !os.IsNotExist(err) {
		err = os.Remove(path)
		if err != nil {
			return errors.Wrap(err, "could not delete the previous active workdir directory symlink")
		}
	}

	return nil
}

// Active returns active working directory
func (h *Handler) Active() (*Workdir, error) {
	path, err := h.activeAbsolutePath()
	if err != nil {
		return nil, err
	}

	resolvedPath, err := filepath.EvalSymlinks(path)
	if os.IsNotExist(err) {
		return nil, ErrMalformed.New("active", err)
	}

	return h.builder.build(resolvedPath)
}

// List returns array of working directories
func (h *Handler) List() ([]*Workdir, error) {
	dirs := make([]string, 0)
	err := filepath.Walk(h.workdirsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		for _, f := range RequiredFiles {
			if !hasContent(path, f) {
				return nil
			}
		}

		dirs = append(dirs, path)
		return nil
	})

	if os.IsNotExist(err) {
		return nil, ErrMalformed.New(h.workdirsPath, err)
	}

	if err != nil {
		return nil, err
	}

	wds := make([]*Workdir, 0, len(dirs))
	for _, p := range dirs {
		wd, err := h.builder.build(p)
		if err != nil {
			return nil, err
		}

		wds = append(wds, wd)
	}

	return wds, nil

}

// Validate validates that the passed working directoy is valid
// It's path must be a directory (or a symlink) containing docker-compose.yml and .env files
func (h *Handler) Validate(w *Workdir) error {
	pointedDir, err := filepath.EvalSymlinks(w.Path)
	if err != nil {
		return ErrMalformed.New(w.Path, "is not a directory")
	}

	if info, err := os.Lstat(pointedDir); err != nil || !info.IsDir() {
		return ErrMalformed.New(pointedDir, "is not a directory")
	}

	for _, f := range RequiredFiles {
		if !hasContent(pointedDir, f) {
			return ErrMalformed.New(pointedDir, fmt.Sprintf("%s not found", f))
		}
	}

	return nil
}

// Remove removes working directory by removing required and optional files,
// and recursively removes directories up to the workdirs root as long as they are empty
func (h *Handler) Remove(w *Workdir) error {
	path := w.Path
	var subPath string
	switch w.Type {
	case Local:
		subPath = "local"
	case Orgs:
		subPath = "orgs"
	}

	basePath := filepath.Join(h.workdirsPath, subPath)

	for _, f := range append(RequiredFiles, OptionalFiles...) {
		file := filepath.Join(path, f)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		if err := os.Remove(file); err != nil {
			return errors.Wrap(err, "could not remove from workdir directory")
		}
	}

	for {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return errors.Wrap(err, "could not read workdir directory")
		}
		if len(files) > 0 {
			return nil
		}

		if err := os.Remove(path); err != nil {
			return errors.Wrap(err, "could not delete workdir directory")
		}

		path = filepath.Dir(path)
		if path == basePath {
			return nil
		}
	}
}

func (h *Handler) activeAbsolutePath() (string, error) {
	return filepath.Join(h.workdirsPath, activeDir), nil
}

package workdir

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/pkg/errors"
	goerrors "gopkg.in/src-d/go-errors.v1"

	composefile "github.com/src-d/sourced-ce/cmd/sourced/compose/file"
	datadir "github.com/src-d/sourced-ce/cmd/sourced/dir"
)

const activeDir = "__active__"

var (
	// RequiredFiles list of required files in a directory to treat it as a working directory
	RequiredFiles = []string{".env", "docker-compose.yml"}

	// OptionalFiles list of optional files that could be deleted when pruning
	OptionalFiles = []string{"docker-compose.override.yml"}

	// ErrMalformed is the returned error when the workdir is wrong
	ErrMalformed = goerrors.NewKind("workdir %s is not valid: %s")
)

// Factory is responsible for the initialization of the workdirs
type Factory struct{}

// InitLocal initializes the workdir for local path and returns the absolute path
func (f *Factory) InitLocal(reposdir string) (string, error) {
	dirName := f.encodeDirName(reposdir)

	envf := envFile{
		Workdir:  dirName,
		ReposDir: reposdir,
	}

	return f.init(dirName, "local", envf)
}

// InitOrgs initializes the workdir for organizationsand returns the absolute path
func (f *Factory) InitOrgs(orgs []string, token string) (string, error) {
	// be indifferent to the order of passed organizations
	sort.Strings(orgs)
	dirName := f.encodeDirName(strings.Join(orgs, ","))

	envf := envFile{
		Workdir:             dirName,
		GithubOrganizations: orgs,
		GithubToken:         token,
	}

	return f.init(dirName, "orgs", envf)
}

func (f *Factory) encodeDirName(dirName string) string {
	return base64.URLEncoding.EncodeToString([]byte(dirName))
}

func (f *Factory) buildAbsPath(dirName, subPath string) (string, error) {
	path, err := workdirsPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, subPath, dirName), nil
}

func (f *Factory) init(dirName string, subPath string, envf envFile) (string, error) {
	workdir, err := f.buildAbsPath(dirName, subPath)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(workdir, 0755)
	if err != nil {
		return "", errors.Wrap(err, "could not create working directory")
	}

	defaultFilePath, err := composefile.InitDefault()
	if err != nil {
		return "", err
	}

	composePath := filepath.Join(workdir, "docker-compose.yml")
	if err := link(defaultFilePath, composePath); err != nil {
		return "", err
	}

	defaultOverridePath, err := composefile.InitDefaultOverride()
	if err != nil {
		return "", err
	}

	workdirOverridePath := filepath.Join(workdir, "docker-compose.override.yml")
	if err := link(defaultOverridePath, workdirOverridePath); err != nil {
		return "", err
	}

	envPath := filepath.Join(workdir, ".env")
	contents := envf.String()
	err = ioutil.WriteFile(envPath, []byte(contents), 0644)

	if err != nil {
		return "", errors.Wrap(err, "could not write .env file")
	}

	return workdir, nil
}

type envFile struct {
	Workdir             string
	ReposDir            string
	GithubOrganizations []string
	GithubToken         string
}

func (f *envFile) String() string {
	volumeType := "bind"
	volumeSource := f.ReposDir
	gitbaseSiva := ""
	if f.ReposDir == "" {
		volumeType = "volume"
		volumeSource = "gitbase_repositories"
		gitbaseSiva = "true"
	}

	return fmt.Sprintf(`COMPOSE_PROJECT_NAME=srcd-%s
	GITBASE_VOLUME_TYPE=%s
	GITBASE_VOLUME_SOURCE=%s
	GITBASE_SIVA=%s
	GITHUB_ORGANIZATIONS=%s
	GITHUB_TOKEN=%s
	`, f.Workdir, volumeType, volumeSource, gitbaseSiva,
		strings.Join(f.GithubOrganizations, ","), f.GithubToken)
}

// SetActive creates a symlink from the fixed active workdir path
// to the workdir for the given repos dir.
func SetActive(workdir string) error {
	activePath, err := absolutePath(activeDir)
	if err != nil {
		return err
	}

	_, err = os.Stat(activePath)
	if !os.IsNotExist(err) {
		err = os.Remove(activePath)
		if err != nil {
			return errors.Wrap(err, "could not delete the previous active workdir directory symlink")
		}
	}

	err = os.Symlink(workdir, activePath)
	if os.IsExist(err) {
		return nil
	}

	return err
}

// UnsetActive removes symlink for active workdir
func UnsetActive() error {
	dir, err := absolutePath(activeDir)
	if err != nil {
		return err
	}

	_, err = os.Lstat(dir)
	if !os.IsNotExist(err) {
		err = os.Remove(dir)
		if err != nil {
			return errors.Wrap(err, "could not delete active workdir directory symlink")
		}
	}

	return nil
}

// Active returns active working directory name
func Active() (string, error) {
	path, err := ActivePath()
	if err != nil {
		return "", err
	}

	wpath, err := workdirsPath()
	if err != nil {
		return "nil", err
	}

	return decodeName(wpath, path)
}

// ActivePath returns absolute path to active working directory
func ActivePath() (string, error) {
	path, err := absolutePath(activeDir)
	if err != nil {
		return "", err
	}

	resolvedPath, err := filepath.EvalSymlinks(path)
	if os.IsNotExist(err) {
		return "", ErrMalformed.New("active", err)
	}

	return resolvedPath, err
}

// List returns array of working directories names
func List() ([]string, error) {
	wpath, err := workdirsPath()
	if err != nil {
		return nil, err
	}

	workdirs, err := ListPaths()
	if err != nil {
		return nil, err
	}

	res := make([]string, len(workdirs))
	for i, d := range workdirs {
		res[i], err = decodeName(wpath, d)
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(res)
	return res, nil
}

// ListPaths returns array of absolute paths to working directories
func ListPaths() ([]string, error) {
	wpath, err := workdirsPath()
	if err != nil {
		return nil, err
	}

	dirs := make(map[string]bool)
	err = filepath.Walk(wpath, func(path string, info os.FileInfo, err error) error {
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

		dirs[path] = true
		return nil
	})

	if os.IsNotExist(err) {
		return nil, ErrMalformed.New(wpath, err)
	}

	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for dir := range dirs {
		res = append(res, dir)
	}

	return res, nil
}

// RemovePath removes working directory by removing required and optional files,
// and recursively removes directories up to the workdirs root as long as they are empty
func RemovePath(path string) error {
	workdirsRoot, err := workdirsPath()
	if err != nil {
		return err
	}

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
		if path == workdirsRoot {
			return nil
		}
	}
}

// SetActivePath similar to SetActive
// but accepts absolute path to a directory instead of a relative one
func SetActivePath(path string) error {
	wpath, err := workdirsPath()
	if err != nil {
		return err
	}

	wd, err := filepath.Rel(wpath, path)
	if err != nil {
		return err
	}

	return SetActive(wd)
}

// ValidatePath validates that the passed dir is valid
// Must be a directory (or a symlink) containing docker-compose.yml and .env files
func ValidatePath(dir string) error {
	pointedDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return ErrMalformed.New(dir, "is not a directory")
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

// path returns the absolute path to
// $HOME/.sourced/workdirs/workdir
func absolutePath(workdir string) (string, error) {
	path, err := workdirsPath()
	if err != nil {
		return "", err
	}

	// On windows replace C:\path with C\path
	workdir = strings.Replace(workdir, ":", "", 1)

	return filepath.Join(path, workdir), nil
}

func hasContent(path, file string) bool {
	empty, err := isEmptyFile(filepath.Join(path, file))
	return !empty && err == nil
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

func link(linkTargetPath, linkPath string) error {
	_, err := os.Stat(linkPath)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return errors.Wrap(err, "could not read the existing FILE_NAME file")
	}

	err = os.Symlink(linkTargetPath, linkPath)
	return errors.Wrap(err, fmt.Sprintf("could not create symlink to %s", linkTargetPath))
}

func workdirsPath() (string, error) {
	path, err := datadir.Path()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "workdirs"), nil
}

// decodeName takes workdirs root and absolute path to workdir
// return human-readable name. It returns an error if the path could not be built
func decodeName(base, target string) (string, error) {
	p, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}

	// workdirs for remote orgs encoded into base64
	decoded, err := base64.URLEncoding.DecodeString(p)
	if err == nil {
		return string(decoded), nil
	}

	// for windows local path convert C\path to C:\path
	if runtime.GOOS == "windows" {
		return string(p[0]) + ":" + p[1:len(p)], nil
	}

	// for *nix prepend root, User/path to /Users/path
	return filepath.Join("/", p), nil
}

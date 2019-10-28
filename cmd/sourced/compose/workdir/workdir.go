package workdir

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	goerrors "gopkg.in/src-d/go-errors.v1"

	datadir "github.com/src-d/sourced-ce/cmd/sourced/dir"
)

const activeDir = "__active__"

var (
	// RequiredFiles list of required files in a directory to treat it as a working directory
	RequiredFiles = []string{".env", "docker-compose.yml"}

	// ErrMalformed is the returned error when the workdir is wrong
	ErrMalformed = goerrors.NewKind("workdir %s is not valid")

	// ErrInitFailed is an error returned on workdir initialization for custom cases
	ErrInitFailed = goerrors.NewKind("initialization failed")
)

// Type defines the type of the workdir
type Type int

const (
	// None refers to a failure in identifying the type of the workdir
	None Type = iota
	// Local refers to a workdir that has been initialized for local repos
	Local
	// Orgs refers to a workdir that has been initialized for organizations
	Orgs
)

// Workdir represents a workdir associated with a local or an orgs initialization
type Workdir struct {
	// Type is the type of working directory
	Type Type
	// Name is a human-friendly string to identify the workdir
	Name string
	// Path is the absolute path corresponding to the workdir
	Path string
}

type builder struct {
	workdirsPath string
}

// build returns the Workdir instance corresponding to the provided absolute path
// the path must be inside `workdirsPath`
func (b *builder) Build(path string) (*Workdir, error) {
	wdType, err := b.typeFromPath(path)
	if err != nil {
		return nil, err
	}

	if wdType == None {
		return nil, fmt.Errorf("invalid workdir type for path %s", path)
	}

	wdName, err := b.workdirName(wdType, path)
	if err != nil {
		return nil, err
	}

	return &Workdir{
		Type: wdType,
		Name: wdName,
		Path: path,
	}, nil
}

// workdirName returns the workdir name given its type and absolute path
func (b *builder) workdirName(wdType Type, path string) (string, error) {
	var subPath string
	switch wdType {
	case Local:
		subPath = "local"
	case Orgs:
		subPath = "orgs"
	}

	encoded, err := filepath.Rel(filepath.Join(b.workdirsPath, subPath), path)
	if err != nil {
		return "", err
	}

	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err == nil {
		return string(decoded), nil
	}

	return "", err
}

// typeFromPath returns the workdir type corresponding to the provided absolute path
func (b *builder) typeFromPath(path string) (Type, error) {
	suffix, err := filepath.Rel(b.workdirsPath, path)
	if err != nil {
		return None, err
	}

	switch filepath.Dir(suffix) {
	case "local":
		return Local, nil
	case "orgs":
		return Orgs, nil
	default:
		return None, nil
	}
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

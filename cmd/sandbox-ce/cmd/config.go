package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-cli.v0"
)

const orgName = "smacker"
const repoName = "superset-compose"
const defaultConfigName = "__default__"
const composerConfigFileURL = "https://raw.githubusercontent.com/%s/%s/%s/docker-compose.yml"

type configCmd struct {
	cli.PlainCommand `name:"config" short-description:"Config"`
}

type configDownloadCmd struct {
	Command `name:"download" short-description:"Download"`

	Revision string `short:"r" long:"revision" description:"revision to download" default:"master"`
}

func (c *configDownloadCmd) Execute(args []string) error {
	confURL := getComposerConfFileURL(c.Revision)
	outPath, err := getConfigPath(c.Revision)
	if err != nil {
		return err
	}

	err = c.downloadConfig(confURL, outPath)
	if err != nil {
		return err
	}

	return nil
}

func (c *configDownloadCmd) downloadConfig(url, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

type configListCmd struct {
	Command `name:"list" short-description:"List"`
}

func (c *configListCmd) Execute(args []string) error {
	confDirPath, err := getConfigDirPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(confDirPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	files, err := ioutil.ReadDir(confDirPath)
	if err != nil {
		return err
	}

	defaultRevision, err := getDefaultConfigRevision()
	if err != nil {
		return err
	}

	for _, f := range files {
		entry := f.Name()
		if entry == defaultConfigName {
			continue
		}

		if entry == defaultRevision {
			entry = fmt.Sprintf("%s (default)", entry)
		}

		fmt.Println(entry)
	}

	return nil
}

type configSetDefaultCmd struct {
	Command `name:"set-default" short-description:"Set default config"`

	Args struct {
		Revision string `positional-arg-name:"revision"`
	} `positional-args:"yes" required:"yes"`
}

func (c *configSetDefaultCmd) Execute(args []string) error {
	revisionConfPath, err := getConfigPath(c.Args.Revision)
	if err != nil {
		return err
	}

	if _, err := os.Stat(revisionConfPath); err != nil {
		if os.IsNotExist(err) {
			return errors.Wrapf(err, "no configuration found for revision `%s`",
				c.Args.Revision)
		}

		return err
	}

	defaultConfPath, err := getConfigPath(defaultConfigName)
	if err != nil {
		return err
	}

	if _, err := os.Lstat(defaultConfPath); err == nil {
		if err := os.Remove(defaultConfPath); err != nil {
			return fmt.Errorf("failed to unlink: %+v", err)
		}
	}

	return os.Symlink(revisionConfPath, defaultConfPath)
}

func getConfigDirPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "could not detect home directory")
	}

	return filepath.Join(homedir, ".srcd", "composer-configs"), nil
}

func getConfigPath(revision string) (string, error) {
	confDirPath, err := getConfigDirPath()
	if err != nil {
		return "", err
	}

	dirPath := filepath.Join(confDirPath, revision)
	if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", errors.Wrapf(err, "error while creating directory %s", dirPath)
	}

	return filepath.Join(dirPath, "docker-compose.yml"), nil
}

func getDefaultConfigRevision() (string, error) {
	defaultConfPath, err := getConfigPath(defaultConfigName)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(defaultConfPath); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", err
	}

	dest, err := filepath.EvalSymlinks(defaultConfPath)
	if err != nil {
		return "", err
	}

	_, rev := filepath.Split(filepath.Dir(dest))
	return rev, nil
}

func getComposerConfFileURL(revision string) string {
	return fmt.Sprintf(composerConfigFileURL, orgName, repoName, revision)
}

func init() {
	c := rootCmd.AddCommand(&configCmd{})
	c.AddCommand(&configDownloadCmd{})
	c.AddCommand(&configListCmd{})
	c.AddCommand(&configSetDefaultCmd{})
}

// Package workdir provides functions to manage docker compose working
// directories inside the $HOME/.sourced/workdirs directory
package workdir

import (
	"crypto/sha1"
	"encoding/hex"
)

// Init creates a working directory in ~/.sourced for the given repositories
// directory. The working directory will contain a docker-compose.yml and a
// .env file.
// If the directory is already initialized the function returns with no error.
// The returned value is the absolute path to $HOME/.sourced/workdirs/reposdir
func InitWithPath(reposdir string) (string, error) {
	workdir, err := absolutePath(reposdir)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum([]byte(reposdir))
	hashSt := hex.EncodeToString(hash[:])
	envf := envFile{
		Workdir:  hashSt,
		ReposDir: reposdir,
	}

	if err := initWorkdir(workdir, envf); err != nil {
		return "", err
	}

	return workdir, nil
}

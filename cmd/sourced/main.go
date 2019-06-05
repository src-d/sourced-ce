package main

import (
	"fmt"
	"runtime"

	"github.com/src-d/sourced-ce/cmd/sourced/cmd"
	composefile "github.com/src-d/sourced-ce/cmd/sourced/compose/file"
	"github.com/src-d/sourced-ce/cmd/sourced/release"
)

// this variable is rewritten during the CI build step
var version = "master"
var build = "dev"

func main() {
	composefile.SetVersion(version)
	cmd.Init(version, build)

	checkUpdates()

	cmd.Execute()
}

func checkUpdates() {
	if version == "master" {
		return
	}

	update, latest, err := release.FindUpdates(version)
	if err != nil {
		return
	}

	if update {
		s := fmt.Sprintf(
			`There is a newer version. Current version: %s, latest version: %s
Please go to https://github.com/src-d/sourced-ce/releases/latest to upgrade.
`, version, latest)

		fmt.Println(yellow(s))
	}
}

func yellow(s string) string {
	if runtime.GOOS == "windows" {
		return s
	}

	return fmt.Sprintf("\x1b[33m%s\x1b[0m", s)
}

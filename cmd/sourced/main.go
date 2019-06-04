package main

import (
	"github.com/src-d/sourced-ce/cmd/sourced/cmd"
	composefile "github.com/src-d/sourced-ce/cmd/sourced/compose/file"
)

// this variable is rewritten during the CI build step
var version = "master"
var build = "dev"

func main() {
	composefile.SetVersion(version)
	cmd.Init(version, build)
	cmd.Execute()
}

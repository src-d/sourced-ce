package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/src-d/sourced-ce/cmd/sourced/compose"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/file"
	"github.com/src-d/sourced-ce/cmd/sourced/compose/workdir"
	"github.com/src-d/sourced-ce/cmd/sourced/dir"
	"github.com/src-d/sourced-ce/cmd/sourced/format"

	"gopkg.in/src-d/go-cli.v0"
)

const name = "sourced"

var version = "master"

var rootCmd = cli.NewNoDefaults(name, "source{d} Community Edition & Enterprise Edition CLI client")

// Init sets the version rewritten by the CI build and adds default sub commands
func Init(v, build string) {
	version = v

	rootCmd.AddCommand(&cli.VersionCommand{
		Name:    name,
		Version: version,
		Build:   build,
	})

	if runtime.GOOS != "windows" {
		rootCmd.AddCommand(&cli.CompletionCommand{
			Name: name,
		}, cli.InitCompletionCommand(name))
	}
}

// Command implements the default group flags. It is meant to be embedded into
// other application commands to provide default behavior for logging, config
type Command struct {
	cli.PlainCommand
	cli.LogOptions `group:"Log Options"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := dir.Prepare(); err != nil {
		fmt.Println(err)
		log(err)
		os.Exit(1)
	}

	if err := rootCmd.Run(os.Args); err != nil {
		log(err)
		os.Exit(1)
	}
}

func log(err error) {
	switch {
	case workdir.ErrMalformed.Is(err) || dir.ErrNotExist.Is(err):
		printRed("Cannot perform this action, source{d} needs to be initialized first with the 'init' sub command")
	case workdir.ErrInitFailed.Is(err):
		printRed("Cannot perform this action, full re-initialization is needed, run 'prune' command first")
	case dir.ErrNotValid.Is(err):
		printRed("Cannot perform this action, config directory is not valid")
	case compose.ErrComposeAlternative.Is(err):
		printRed("docker-compose is not installed, and there was an error while trying to use the container alternative")
		printRed("  see: https://docs.sourced.tech/community-edition/quickstart/1-install-requirements#docker-compose")
	case file.ErrConfigDownload.Is(err):
		printRed("The source{d} CE config file could not be set as active")
	case fmt.Sprintf("%T", err) == "*flags.Error":
		// syntax error is already logged by go-cli
	default:
		// unknown errors have no special message
	}

	switch {
	case dir.ErrNetwork.Is(err):
		// TODO(dpordomingo): if start using "https://golang.org/pkg/errors/",
		//					  we could do `var myErr ErrNetwork; errors.As(err, &myErr)`
		// 					  to provide more info about the actual network error
		printRed("The resource could not be downloaded from the Internet")
		printRed("  see: https://docs.sourced.tech/community-edition/quickstart/1-install-requirements#internet-connection")
	case dir.ErrWrite.Is(err):
		// TODO(dpordomingo): see todo above
		printRed("could not write file")

	}
}

func printRed(message string) {
	fmt.Println(format.Colorize(format.Red, message))
}

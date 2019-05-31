package cmd

import (
	"gopkg.in/src-d/go-cli.v0"
)

const name = "sourced"

var rootCmd = cli.NewNoDefaults(name, "source{d} Sandbox")

// Command implements the default group flags. It is meant to be embedded into
// other application commands to provide default behavior for logging, config
type Command struct {
	cli.PlainCommand
	cli.LogOptions `group:"Log Options"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.RunMain()
}

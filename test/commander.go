// +build integration

package test

import (
	"fmt"
	"os"
	"time"

	"gotest.tools/icmd"
)

type Commander struct {
	bin        string
	sourcedDir string
}

// RunCmd runs a command with the appropriate SOURCED_DIR and returns a Result
func (s *Commander) RunCmd(args []string, cmdOperators ...icmd.CmdOp) *icmd.Result {
	c := icmd.Command(s.bin, args...)

	newEnv := append(os.Environ(),
		fmt.Sprintf("SOURCED_DIR=%s", s.sourcedDir))
	op := append(cmdOperators, icmd.WithEnv(newEnv...))

	return icmd.RunCmd(c, op...)
}

// RunCommand is a convenience wrapper for RunCmd that accepts variadic arguments.
// It runs a command with the appropriate SOURCED_DIR and returns a Result
func (s *Commander) RunCommand(args ...string) *icmd.Result {
	return s.RunCmd(args)
}

// RunCommandWithTimeout runs a command with the given timeout
func (s *Commander) RunCommandWithTimeout(timeout time.Duration, args ...string) *icmd.Result {
	return s.RunCmd(args, icmd.WithTimeout(timeout))
}

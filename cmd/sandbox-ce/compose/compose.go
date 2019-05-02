package compose

import (
	"context"
	"os"
	"os/exec"
)

func Run(ctx context.Context, arg ...string) error {
	cmd := exec.CommandContext(ctx, "docker-compose", arg...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

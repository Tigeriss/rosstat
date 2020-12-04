package debug

import (
	"context"
	"os"
	"os/exec"
)

func RunReactDevServer(ctx context.Context) {
	cmd := exec.CommandContext(ctx, "yarn", "start")
	cmd.Env = append(cmd.Env, "BROWSER=none")
	cmd.Dir = "ui"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}
}

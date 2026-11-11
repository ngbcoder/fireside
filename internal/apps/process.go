package apps

import (
	"context"
	"os/exec"
)

// cmdProcess wraps an *exec.Cmd as a Process.
type cmdProcess struct {
	cmd *exec.Cmd
}

func (p *cmdProcess) Wait() error { return p.cmd.Wait() }

func (p *cmdProcess) Stop() error {
	if p.cmd.Process == nil {
		return nil
	}
	return p.cmd.Process.Kill()
}

// runCommand builds and starts a command, returning a Process.
func runCommand(ctx context.Context, name string, args ...string) (Process, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &cmdProcess{cmd: cmd}, nil
}

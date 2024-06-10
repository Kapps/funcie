package tools

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// ProcessRunner provides a helper interface for running shell commands.
type ProcessRunner interface {
	// Run runs a shell command with the specified arguments and returns the output once the command completes.
	Run(cmd string, baseDir string, args ...string) (string, error)
}

type processRunner struct{}

// NewProcessRunner creates a new ProcessRunner.
func NewProcessRunner() ProcessRunner {
	return &processRunner{}
}

func (p *processRunner) Run(cmd string, baseDir string, args ...string) (string, error) {
	proc := exec.Command(cmd, args...)

	var out bytes.Buffer
	proc.Stdout = io.MultiWriter(os.Stdout, &out)
	proc.Stderr = os.Stderr

	proc.Dir = baseDir

	// Ensure process is killed if the parent process is killed
	// TODO: Validate this and if other signals are needed
	proc.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := proc.Run(); err != nil {
		return "", fmt.Errorf("failed to run command %v: %w", cmd, err)
	}

	return out.String(), nil
}

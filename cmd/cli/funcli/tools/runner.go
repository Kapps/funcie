package tools

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunnerOpts provides options for running shell commands.
type RunnerOpts struct {
	// Args are the arguments to pass to the command.
	Args []string
	// Env is the environment variables to set for the command, in the form of key-value mappings.
	Env map[string]string
	// WorkingDir is the working directory to run the command in, or an empty string to use the current working directory.
	WorkingDir string
}

// ProcessRunner provides a helper interface for running shell commands.
type ProcessRunner interface {
	// Run runs a shell command with the specified arguments and returns the output once the command completes.
	Run(cmd string, args ...string) (string, error)
	// RunWithOptions runs a shell command with the specified options and returns the output once the command completes.
	RunWithOptions(cmd string, options RunnerOpts) (string, error)
}

type processRunner struct{}

// NewProcessRunner creates a new ProcessRunner.
func NewProcessRunner() ProcessRunner {
	return &processRunner{}
}

func (p *processRunner) Run(cmd string, args ...string) (string, error) {
	return p.RunWithOptions(cmd, RunnerOpts{
		Args: args,
	})
}

func (p *processRunner) RunWithOptions(cmd string, options RunnerOpts) (string, error) {
	proc := exec.Command(cmd, options.Args...)

	var out bytes.Buffer
	proc.Stdout = io.MultiWriter(os.Stdout, &out)
	proc.Stderr = os.Stderr

	if options.Env != nil {
		for k, v := range options.Env {
			proc.Env = append(proc.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	proc.Env = append(proc.Env, os.Environ()...)

	if options.WorkingDir != "" {
		proc.Dir = options.WorkingDir
	}

	// Ensure process is killed if the parent process is killed
	// TODO: Validate this and if other signals are needed
	proc.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	fmt.Println(color.HiBlackString(fmt.Sprintf("%v %v", cmd, strings.Join(options.Args, " "))))

	if err := proc.Run(); err != nil {
		return "", fmt.Errorf("failed to run command %v: %w", cmd, err)
	}

	return out.String(), nil
}

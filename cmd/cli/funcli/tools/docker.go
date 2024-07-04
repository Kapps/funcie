package tools

import "fmt"

// DockerRunOptions provides options for how to run a container.
type DockerRunOptions struct {
	// Env is a map of environment variables to set in the container.
	Env map[string]string
	// ExposedPorts is a map of ports to expose on the container.
	ExposedPorts []int
	// RestartPolicy is the policy to use when the container exits.
	RestartPolicy string
}

// DockerClient provides an interface for interacting with Docker containers.
type DockerClient interface {
	// RunContainer runs a container with the specified image and options.
	RunContainer(image string, opts DockerRunOptions) error
}

type dockerCliClient struct {
	runner ProcessRunner
}

// NewDockerCliClient creates a new DockerClient that wraps the Docker CLI tool.
func NewDockerCliClient(runner ProcessRunner) DockerClient {
	return &dockerCliClient{
		runner: runner,
	}
}

func (d *dockerCliClient) RunContainer(image string, opts DockerRunOptions) error {
	args := []string{"run", "-d"}

	for k, v := range opts.Env {
		args = append(args, "-e", k+"="+v)
	}

	for _, port := range opts.ExposedPorts {
		args = append(args, "-p", fmt.Sprintf("%d", port))
	}

	if opts.RestartPolicy != "" {
		args = append(args, "--restart", opts.RestartPolicy, image)
	}

	_, err := d.runner.Run("docker", args...)
	if err != nil {
		return fmt.Errorf("failed to run container %v: %w", image, err)
	}

	return nil
}

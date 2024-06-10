package tools

import (
	"fmt"
	"os"
)

// GitClient provides an interface for interacting with Git repositories.
type GitClient interface {
	// ShallowClone clones a repository to the specified directory with a depth of 1.
	// If the directory already exists, it will be deleted and recreated with the new repository.
	ShallowClone(url string, directory string, branch string) error
}

type gitCliClient struct {
	runner ProcessRunner
}

// We could use a pure go implementation, but that would bring in a lot of dependencies.
// Instead, we'll use the git CLI tool to clone repositories.

// NewGitCliClient creates a new GitClient that wraps the Git CLI tool.
func NewGitCliClient(runner ProcessRunner) GitClient {
	return &gitCliClient{
		runner: runner,
	}
}

func (g *gitCliClient) ShallowClone(url string, directory string, branch string) error {
	if err := os.RemoveAll(directory); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove directory %v: %w", directory, err)
	}

	_, err := g.runner.Run("git", "", "clone", url, "--depth", "1", directory)
	if err != nil {
		return fmt.Errorf("failed to clone repository %v: %w", url, err)
	}

	return nil
}

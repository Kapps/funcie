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
	// Checkout checks out a branch in the specified directory.
	// If the directory does not exist, it will be created and the repository cloned.
	// If the directory exists, it will be updated with the latest changes on the given branch from the repository.
	Checkout(url string, directory string, branch string) error
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

	_, err := g.runner.Run("git", "clone", url, "-c", "advice.detachedHead=false", "--branch", branch, "--depth", "1", directory)
	if err != nil {
		return fmt.Errorf("failed to clone repository %v: %w", url, err)
	}

	return nil
}

func (g *gitCliClient) Checkout(url string, directory string, branch string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return g.ShallowClone(url, directory, branch)
	}

	_, err := g.runner.Run("git", "-C", directory, "fetch", "--tags", "origin")
	if err != nil {
		return fmt.Errorf("failed to fetch branch %v: %w", branch, err)
	}

	_, err = g.runner.Run("git", "-c", "advice.detachedHead=false", "-C", directory, "checkout", branch)
	if err != nil {
		return fmt.Errorf("failed to checkout branch %v: %w", branch, err)
	}

	_, err = g.runner.Run("git", "-C", directory, "pull")
	if err != nil {
		return fmt.Errorf("failed to pull changes for branch %v: %w", branch, err)
	}

	return nil
}

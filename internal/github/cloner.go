package github

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Cloner handles git clone operations with safety measures.
type Cloner struct {
	// Token is an optional GitHub personal access token for private repos.
	// If empty, will try to use GITHUB_TOKEN environment variable.
	Token string

	// Shallow enables shallow cloning (--depth=1) for faster operations.
	Shallow bool
}

// NewCloner creates a new Cloner with default settings.
func NewCloner() *Cloner {
	return &Cloner{
		Shallow: true, // Default to shallow clone for performance
	}
}

// Clone clones a GitHub repository to a temporary directory.
// Returns the CloneResult with the local path and commit information.
// The caller is responsible for cleaning up the temporary directory.
func (c *Cloner) Clone(info *RepoInfo) (*CloneResult, error) {
	// Validate repository info
	if err := ValidateRepoInfo(info); err != nil {
		return nil, fmt.Errorf("invalid repository: %w", err)
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git is not installed or not in PATH: %w", err)
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "rpg-github-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Build clone URL with authentication if token provided
	cloneURL := c.buildAuthenticatedURL(info.CloneURL)

	// Build git clone command
	args := []string{"clone"}

	// Safety: Disable git hooks during clone
	args = append(args, "--config", "core.hooksPath=/dev/null")

	// Shallow clone if enabled
	if c.Shallow {
		args = append(args, "--depth=1")
	}

	// Add branch/tag/ref if specified
	if info.Ref != "" {
		args = append(args, "--branch", info.Ref)
	}

	// Target directory within temp
	targetDir := filepath.Join(tempDir, info.Name)
	args = append(args, cloneURL, targetDir)

	// Execute git clone
	cmd := exec.Command("git", args...)
	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0", // Disable interactive prompts
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Clean up temp directory on failure
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("git clone failed: %w\nOutput: %s", err, string(output))
	}

	// Get the current branch
	branch, err := c.getBranch(targetDir)
	if err != nil {
		branch = "unknown"
	}

	// Get the HEAD commit SHA
	commitSHA, err := c.getCommitSHA(targetDir)
	if err != nil {
		commitSHA = "unknown"
	}

	return &CloneResult{
		LocalPath: targetDir,
		Branch:    branch,
		CommitSHA: commitSHA,
	}, nil
}

// buildAuthenticatedURL adds authentication to the clone URL if a token is available.
func (c *Cloner) buildAuthenticatedURL(url string) string {
	token := c.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		return url
	}

	// Convert https://github.com/owner/repo.git to https://token@github.com/owner/repo.git
	if strings.HasPrefix(url, "https://github.com/") {
		return strings.Replace(url, "https://github.com/", fmt.Sprintf("https://%s@github.com/", token), 1)
	}
	return url
}

// getBranch returns the current branch name.
func (c *Cloner) getBranch(repoDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getCommitSHA returns the current HEAD commit SHA.
func (c *Cloner) getCommitSHA(repoDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Cleanup removes a cloned repository directory.
func Cleanup(path string) error {
	if path == "" {
		return nil
	}
	// Safety check: only remove directories under system temp
	tempDir := os.TempDir()
	if !strings.HasPrefix(path, tempDir) {
		return fmt.Errorf("refusing to remove directory outside temp: %s", path)
	}
	return os.RemoveAll(path)
}

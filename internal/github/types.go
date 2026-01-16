// Package github provides GitHub repository cloning and analysis functionality.
package github

// RepoInfo contains parsed information about a GitHub repository.
type RepoInfo struct {
	Owner    string // Repository owner (user or organization)
	Name     string // Repository name
	Ref      string // Branch, tag, or commit (optional)
	CloneURL string // Full HTTPS clone URL
}

// CloneResult contains the result of a clone operation.
type CloneResult struct {
	LocalPath string // Path to cloned repository
	Branch    string // Branch that was checked out
	CommitSHA string // HEAD commit SHA
}

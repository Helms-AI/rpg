package github

import (
	"fmt"
	"regexp"
	"strings"
)

// Regex patterns for different GitHub URL formats
var (
	// https://github.com/owner/repo or https://github.com/owner/repo.git
	httpsPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+?)(?:\.git)?$`)

	// github.com/owner/repo
	noProtocolPattern = regexp.MustCompile(`^github\.com/([^/]+)/([^/]+?)(?:\.git)?$`)

	// git@github.com:owner/repo.git
	sshPattern = regexp.MustCompile(`^git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`)

	// owner/repo (shorthand)
	shorthandPattern = regexp.MustCompile(`^([a-zA-Z0-9][-a-zA-Z0-9]*)/([a-zA-Z0-9._-]+)$`)
)

// ParseRepository parses a GitHub repository URL or shorthand into RepoInfo.
// Supported formats:
//   - https://github.com/owner/repo
//   - https://github.com/owner/repo.git
//   - github.com/owner/repo
//   - git@github.com:owner/repo.git
//   - owner/repo (shorthand)
//   - owner/repo@branch (with ref)
//   - owner/repo#commit (with commit SHA)
func ParseRepository(input string) (*RepoInfo, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("repository cannot be empty")
	}

	// Extract ref if present (owner/repo@branch or owner/repo#commit)
	var ref string
	if idx := strings.LastIndex(input, "@"); idx > 0 && !strings.Contains(input, "://") && !strings.HasPrefix(input, "git@") {
		ref = input[idx+1:]
		input = input[:idx]
	} else if idx := strings.LastIndex(input, "#"); idx > 0 {
		ref = input[idx+1:]
		input = input[:idx]
	}

	var owner, name string

	// Try each pattern
	if matches := httpsPattern.FindStringSubmatch(input); len(matches) == 3 {
		owner, name = matches[1], matches[2]
	} else if matches := noProtocolPattern.FindStringSubmatch(input); len(matches) == 3 {
		owner, name = matches[1], matches[2]
	} else if matches := sshPattern.FindStringSubmatch(input); len(matches) == 3 {
		owner, name = matches[1], matches[2]
	} else if matches := shorthandPattern.FindStringSubmatch(input); len(matches) == 3 {
		owner, name = matches[1], matches[2]
	} else {
		return nil, fmt.Errorf("invalid GitHub repository format: %s\nSupported formats: owner/repo, https://github.com/owner/repo, git@github.com:owner/repo.git", input)
	}

	// Clean up name (remove .git suffix if present)
	name = strings.TrimSuffix(name, ".git")

	return &RepoInfo{
		Owner:    owner,
		Name:     name,
		Ref:      ref,
		CloneURL: fmt.Sprintf("https://github.com/%s/%s.git", owner, name),
	}, nil
}

// ValidateRepoInfo validates the parsed repository information.
func ValidateRepoInfo(info *RepoInfo) error {
	if info.Owner == "" {
		return fmt.Errorf("repository owner cannot be empty")
	}
	if info.Name == "" {
		return fmt.Errorf("repository name cannot be empty")
	}
	// Basic validation of owner/name format
	if strings.ContainsAny(info.Owner, " \t\n") {
		return fmt.Errorf("invalid owner: %s", info.Owner)
	}
	if strings.ContainsAny(info.Name, " \t\n") {
		return fmt.Errorf("invalid repository name: %s", info.Name)
	}
	return nil
}

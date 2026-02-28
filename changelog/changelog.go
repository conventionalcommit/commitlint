package changelog

import (
	"time"

	"github.com/conventionalcommit/commitlint/commit"
)

// CommitInfo holds a parsed commit with git metadata
type CommitInfo struct {
	// Commit is the parsed conventional commit
	Commit commit.Commit

	// Hash is the full commit SHA
	Hash string

	// ShortHash is the abbreviated commit SHA (7 chars)
	ShortHash string

	// CommitURL is the full URL to the commit
	CommitURL string

	// Author is the commit author name
	Author string

	// Date is the commit date
	Date time.Time

	// References holds extracted issue references
	References []string
}

// TypeGroup holds commits of one type
type TypeGroup struct {
	// Type is the conventional commit type
	Type string

	// Header is the display header for this group
	Header string

	// Commits holds all commits of this type
	Commits []CommitInfo
}

// VersionChangelog holds all commits for one version/tag range
type VersionChangelog struct {
	// Version is the tag name (e.g., "v1.0.0") or "Unreleased"
	Version string

	// Date is the release date formatted as YYYY-MM-DD
	Date string

	// CompareURL is the URL to compare with previous version
	CompareURL string

	// FromRef is the start reference (exclusive)
	FromRef string

	// ToRef is the end reference (inclusive)
	ToRef string

	// Groups holds commits grouped by type
	Groups []TypeGroup

	// Breaking holds breaking change commits
	Breaking []CommitInfo

	// Other holds non-conventional commits
	Other []CommitInfo
}

// Changelog is the full changelog document
type Changelog struct {
	// Header is the top-level changelog header
	Header string

	// Versions holds all version changelogs, newest first
	Versions []VersionChangelog
}

// Formatter formats a Changelog into string output
type Formatter interface {
	// Name returns the name of the formatter
	Name() string

	// Format formats the full changelog
	Format(changelog *Changelog) (string, error)

	// FormatVersion formats a single version changelog
	FormatVersion(version *VersionChangelog, cfg *Config) (string, error)
}

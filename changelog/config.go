// Package changelog contains core types for changelog generation
package changelog

// Config holds the changelog generation configuration
type Config struct {
	// Formatter is the name of the formatter to use
	Formatter string `yaml:"formatter"`

	// Output is the output file path, "-" for stdout
	Output string `yaml:"output"`

	// Header is the top-level changelog header
	Header string `yaml:"header"`

	// Repository holds repository URL configuration
	Repository RepositoryConfig `yaml:"repository"`

	// IssuePrefixes defines prefixes for issue references
	IssuePrefixes []string `yaml:"issue-prefixes"`

	// IncludeOther includes non-conventional commits
	IncludeOther bool `yaml:"include-other"`

	// IncludeBreaking adds a separate breaking changes section
	IncludeBreaking bool `yaml:"include-breaking"`

	// SkipMergeCommits skips merge commits from changelog
	SkipMergeCommits bool `yaml:"skip-merge-commits"`

	// Types defines commit type grouping and display
	Types []TypeConfig `yaml:"types"`
}

// RepositoryConfig holds repository URL configuration
type RepositoryConfig struct {
	// URL is the base repository URL (auto-inferred from git remote)
	URL string `yaml:"url"`

	// CommitURL is the template for commit links
	// Placeholders: {{hash}}
	CommitURL string `yaml:"commit-url"`

	// CompareURL is the template for version compare links
	// Placeholders: {{from}}, {{to}}
	CompareURL string `yaml:"compare-url"`
}

// TypeConfig defines how a commit type is displayed in the changelog
type TypeConfig struct {
	// Type is the conventional commit type (e.g., "feat", "fix")
	Type string `yaml:"type"`

	// Header is the display header (e.g., "Features", "Bug Fixes")
	Header string `yaml:"header"`

	// Hidden if true, commits of this type are excluded from output
	Hidden bool `yaml:"hidden"`
}

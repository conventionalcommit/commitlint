package config

import "github.com/conventionalcommit/commitlint/lint"

// LintMessage lints commitMsg using the default configuration.
// It is the simplest entry point for programmatic use: no config file is needed.
//
// For custom configuration use Parse or NewDefault, then NewLinter.
func LintMessage(commitMsg string) (*lint.Result, error) {
	linter, err := NewLinter(NewDefault())
	if err != nil {
		return nil, err
	}
	return linter.ParseAndLint(commitMsg)
}

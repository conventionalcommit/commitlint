package config

import (
	"github.com/conventionalcommit/commitlint/changelog"
	"github.com/conventionalcommit/commitlint/lint"
)

// LintMessage lints commitMsg using the default configuration.
// It is the simplest entry point for programmatic use: no config file is needed.
//
// For custom configuration use Parse or NewDefault, then NewLinter.
func LintMessage(commitMsg string) (*lint.Result, error) {
	conf := NewDefaultLint()
	linter, err := NewLinter(conf)
	if err != nil {
		return nil, err
	}
	return linter.ParseAndLint(commitMsg)
}

// GenerateChangelog generates changelog using the default configuration.
// It is the simplest entry point for programmatic use.
func GenerateChangelog(repoDir string) (*changelog.Changelog, error) {
	clConf := NewDefaultChangelog()
	gen, err := changelog.New(repoDir, clConf)
	if err != nil {
		return nil, err
	}
	return gen.GenerateAll()
}

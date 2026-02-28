package config

import (
	"github.com/conventionalcommit/commitlint/changelog"
	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/internal/casing"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/lint/formatter"
	"github.com/conventionalcommit/commitlint/lint/rule"
)

const (
	DefaultTypeCharset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultScopeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/,"
	DefaultTrailer      = "Signed-off-by"
)

// Rule name variables — avoids repeated (&rule.XxxRule{}).Name() allocations.
var (
	ruleHeadMinLen     = (&rule.HeadMinLenRule{}).Name()
	ruleHeadMaxLen     = (&rule.HeadMaxLenRule{}).Name()
	ruleBodyMinLen     = (&rule.BodyMinLenRule{}).Name()
	ruleBodyMaxLen     = (&rule.BodyMaxLenRule{}).Name()
	ruleBodyMaxLineLen = (&rule.BodyMaxLineLenRule{}).Name()
	ruleFooterMinLen   = (&rule.FooterMinLenRule{}).Name()
	ruleFooterMaxLen   = (&rule.FooterMaxLenRule{}).Name()
	ruleFooterMaxLine  = (&rule.FooterMaxLineLenRule{}).Name()
	ruleTypeMinLen     = (&rule.TypeMinLenRule{}).Name()
	ruleTypeMaxLen     = (&rule.TypeMaxLenRule{}).Name()
	ruleScopeMinLen    = (&rule.ScopeMinLenRule{}).Name()
	ruleScopeMaxLen    = (&rule.ScopeMaxLenRule{}).Name()
	ruleDescMinLen     = (&rule.DescriptionMinLenRule{}).Name()
	ruleDescMaxLen     = (&rule.DescriptionMaxLenRule{}).Name()

	ruleTypeEnum       = (&rule.TypeEnumRule{}).Name()
	ruleScopeEnum      = (&rule.ScopeEnumRule{}).Name()
	ruleFooterEnum     = (&rule.FooterEnumRule{}).Name()
	ruleFooterTypeEnum = (&rule.FooterTypeEnumRule{}).Name()

	ruleTypeCharset  = (&rule.TypeCharsetRule{}).Name()
	ruleScopeCharset = (&rule.ScopeCharsetRule{}).Name()

	ruleTypeCase  = (&rule.TypeCaseRule{}).Name()
	ruleScopeCase = (&rule.ScopeCaseRule{}).Name()
	ruleDescCase  = (&rule.DescriptionCaseRule{}).Name()
	ruleBodyCase  = (&rule.BodyCaseRule{}).Name()
	ruleHeadCase  = (&rule.HeaderCaseRule{}).Name()

	ruleTypeEmpty   = (&rule.TypeEmptyRule{}).Name()
	ruleScopeEmpty  = (&rule.ScopeEmptyRule{}).Name()
	ruleBodyEmpty   = (&rule.BodyEmptyRule{}).Name()
	ruleFooterEmpty = (&rule.FooterEmptyRule{}).Name()
	ruleDescEmpty   = (&rule.DescriptionEmptyRule{}).Name()

	ruleHeadFullStop = (&rule.HeaderFullStopRule{}).Name()
	ruleBodyFullStop = (&rule.BodyFullStopRule{}).Name()
	ruleDescFullStop = (&rule.DescriptionFullStopRule{}).Name()

	ruleBodyLeadingBlank   = (&rule.BodyLeadingBlankRule{}).Name()
	ruleFooterLeadingBlank = (&rule.FooterLeadingBlankRule{}).Name()

	ruleHeaderTrim = (&rule.HeaderTrimRule{}).Name()

	ruleSignedOffBy   = (&rule.SignedOffByRule{}).Name()
	ruleTrailerExists = (&rule.TrailerExistsRule{}).Name()

	ruleBreakingExcl = (&rule.BreakingChangeExclamationMarkRule{}).Name()

	defaultLintFormatter = (&formatter.DefaultFormatter{}).Name()
)

// NewDefault returns the default Config with lint and changelog defaults.
func NewDefault() *Config {
	return &Config{
		Lint:      NewDefaultLint(),
		Changelog: NewDefaultChangelog(),
	}
}

// NewDefaultLint returns the default lint configuration.
func NewDefaultLint() *lint.Config {
	return &lint.Config{
		MinVersion: internal.Version(),
		Formatter:  defaultLintFormatter,
		Rules: []string{
			ruleHeadMinLen,
			ruleHeadMaxLen,
			ruleBodyMaxLineLen,
			ruleFooterMaxLine,
			ruleTypeEnum,
		},
		Severity: lint.SeverityConfig{
			Default: lint.SeverityError,
		},
		Settings: map[string]lint.RuleSetting{
			// Length rules
			ruleHeadMinLen:     {Argument: 10},
			ruleHeadMaxLen:     {Argument: 72},
			ruleBodyMinLen:     {Argument: 0},
			ruleBodyMaxLen:     {Argument: -1},
			ruleBodyMaxLineLen: {Argument: 100},
			ruleFooterMinLen:   {Argument: 0},
			ruleFooterMaxLen:   {Argument: -1},
			ruleFooterMaxLine:  {Argument: 100},
			ruleTypeMinLen:     {Argument: 0},
			ruleTypeMaxLen:     {Argument: -1},
			ruleScopeMinLen:    {Argument: 0},
			ruleScopeMaxLen:    {Argument: -1},
			ruleDescMinLen:     {Argument: 0},
			ruleDescMaxLen:     {Argument: -1},

			// Enum rules
			ruleTypeEnum: {Argument: DefaultTypeEnums()},
			ruleScopeEnum: {
				Argument: []interface{}{},
				Flags:    map[string]interface{}{"allow-empty": true},
			},
			ruleFooterEnum:     {Argument: []interface{}{}},
			ruleFooterTypeEnum: {Argument: []interface{}{}},

			// Charset rules
			ruleTypeCharset:  {Argument: DefaultTypeCharset},
			ruleScopeCharset: {Argument: DefaultScopeCharset},

			// Case rules
			ruleTypeCase:  {Argument: casing.Lower},
			ruleScopeCase: {Argument: casing.Lower},
			ruleDescCase:  {Argument: casing.Lower},
			ruleBodyCase:  {Argument: casing.Lower},
			ruleHeadCase:  {Argument: casing.Lower},

			// Full-stop rules
			ruleHeadFullStop: {Argument: "."},
			ruleBodyFullStop: {Argument: "."},
			ruleDescFullStop: {Argument: "."},

			// Trailer / Signed-off-by
			ruleSignedOffBy:   {Argument: DefaultTrailer},
			ruleTrailerExists: {Argument: DefaultTrailer},

			// Empty rules
			ruleTypeEmpty:   {},
			ruleScopeEmpty:  {},
			ruleBodyEmpty:   {},
			ruleFooterEmpty: {},
			ruleDescEmpty:   {},

			// Leading-blank rules
			ruleBodyLeadingBlank:   {},
			ruleFooterLeadingBlank: {},

			// Header trim
			ruleHeaderTrim: {},

			// Breaking change
			ruleBreakingExcl: {},
		},
		DefaultIgnorePatterns: DefaultIgnorePatterns(),
	}
}

// NewDefaultChangelog returns the default changelog configuration
func NewDefaultChangelog() *changelog.Config {
	return &changelog.Config{
		Formatter:        "markdown",
		Output:           "",
		Header:           "# Changelog",
		IssuePrefixes:    DefaultIssuePrefixes(),
		IncludeOther:     false,
		IncludeBreaking:  true,
		SkipMergeCommits: true,
		Types:            DefaultChangelogTypes(),
	}
}

// DefaultIgnorePatterns returns the default list of ignore patterns
// These patterns match commit messages auto-generated by git commands
// like merge, revert, fixup, squash, etc.
func DefaultIgnorePatterns() []string {
	return []string{
		// GitHub / GitLab merge
		`^Merge pull request #\d+`,
		`^Merge .+ into .+`,
		`^Merge branch '.+'`,
		`^Merge tag '.+'`,
		`^Merge remote-tracking branch '.+'`,

		// Azure DevOps / Bitbucket merge
		`^Merged .+ (in|into) .+`,
		`^Merged PR #?\d+`,

		// Revert and Reapply
		`^(R|r)evert `,
		`^(R|r)eapply `,

		// Fixup, Amend, Squash (git commit --fixup/--squash)
		`^(amend|fixup|squash)! `,

		// Automatic merges
		`^Automatic merge`,
		`^Auto-merged .+ into .+`,

		// Initial commit
		`^Initial commit$`,
	}
}

// DefaultTypeEnums returns the default list of type enums
func DefaultTypeEnums() []interface{} {
	return []interface{}{
		"feat", "fix", "docs", "style", "refactor", "perf",
		"test", "build", "ci", "chore", "revert",
	}
}

// DefaultChangelogTypes returns the default conventional commit type configuration for changelog
func DefaultChangelogTypes() []changelog.TypeConfig {
	return []changelog.TypeConfig{
		{Type: "feat", Header: "Features"},
		{Type: "fix", Header: "Bug Fixes"},
		{Type: "docs", Header: "Documentation"},
		{Type: "style", Header: "Styles", Hidden: true},
		{Type: "refactor", Header: "Refactor", Hidden: true},
		{Type: "perf", Header: "Performance"},
		{Type: "test", Header: "Tests", Hidden: true},
		{Type: "build", Header: "Build", Hidden: true},
		{Type: "ci", Header: "CI", Hidden: true},
		{Type: "chore", Header: "Chores", Hidden: true},
		{Type: "revert", Header: "Reverts", Hidden: true},
	}
}

// DefaultIssuePrefixes returns the default issue prefixes
func DefaultIssuePrefixes() []string {
	return []string{"#"}
}

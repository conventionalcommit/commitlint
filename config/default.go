package config

import (
	"github.com/conventionalcommit/commitlint/formatter"
	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/rule"
)

// Default returns default config
func Default() *lint.Config {
	def := &lint.Config{
		MinVersion: internal.Version(),

		Formatter: (&formatter.DefaultFormatter{}).Name(),

		Rules: map[string]lint.RuleConfig{
			// Header Min Len Rule
			(&rule.HeadMinLenRule{}).Name(): {
				Enabled:  true,
				Severity: lint.SeverityError,
				Argument: 10,
			},

			// Header Max Len Rule
			(&rule.HeadMaxLenRule{}).Name(): {
				Enabled:  true,
				Severity: lint.SeverityError,
				Argument: 50,
			},

			// Body Max Line Rule
			(&rule.BodyMaxLineLenRule{}).Name(): {
				Enabled:  true,
				Severity: lint.SeverityError,
				Argument: 72,
			},

			// Footer Max Line Rule
			(&rule.FooterMaxLineLenRule{}).Name(): {
				Enabled:  true,
				Severity: lint.SeverityError,
				Argument: 72,
			},

			// Types Enum Rule
			(&rule.TypeEnumRule{}).Name(): {
				Enabled:  true,
				Severity: lint.SeverityError,
				Argument: []interface{}{
					"feat", "fix", "docs", "style", "refactor", "perf",
					"test", "build", "ci", "chore", "revert", "merge",
				},
			},

			// Scope Enum Rule
			(&rule.ScopeEnumRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: []interface{}{},
				Flags: map[string]interface{}{
					"allow-empty": true,
				},
			},

			// Body Min Len Rule
			(&rule.BodyMinLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: 0,
			},

			// Body Max Len Rule
			(&rule.BodyMaxLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: -1,
			},

			// Footer Min Len Rule
			(&rule.FooterMinLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: 0,
			},

			// Footer Max Len Rule
			(&rule.FooterMaxLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: -1,
			},

			// Type Min Len Rule
			(&rule.TypeMinLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: 0,
			},

			// Type Max Len Rule
			(&rule.TypeMaxLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: -1,
			},

			// Scope Min Len Rule
			(&rule.ScopeMinLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: 0,
			},

			// Scope Max Len Rule
			(&rule.ScopeMaxLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: -1,
			},

			// Description Min Len Rule
			(&rule.DescriptionMinLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: 0,
			},

			// Description Max Len Rule
			(&rule.DescriptionMaxLenRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: -1,
			},

			// Type Charset Rule
			(&rule.TypeCharsetRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			},

			// Scope Charset Rule
			(&rule.ScopeCharsetRule{}).Name(): {
				Enabled:  false,
				Severity: lint.SeverityError,
				Argument: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/,",
			},
		},
	}
	return def
}

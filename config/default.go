package config

import (
	"github.com/conventionalcommit/commitlint/formatter"
	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/rule"
)

// Default returns default config
func Default() *lint.Config {
	// Enabled Rules
	rules := []string{
		(&rule.HeadMinLenRule{}).Name(),
		(&rule.HeadMaxLenRule{}).Name(),
		(&rule.BodyMaxLineLenRule{}).Name(),
		(&rule.FooterMaxLineLenRule{}).Name(),
		(&rule.TypeEnumRule{}).Name(),
	}

	// Severity Levels
	severity := lint.SeverityConfig{
		Default: lint.SeverityError,
	}

	// Default Rule Settings
	settings := map[string]lint.RuleSetting{
		// Header Min Len Rule
		(&rule.HeadMinLenRule{}).Name(): {
			Argument: 10,
		},

		// Header Max Len Rule
		(&rule.HeadMaxLenRule{}).Name(): {
			Argument: 50,
		},

		// Body Max Line Rule
		(&rule.BodyMaxLineLenRule{}).Name(): {
			Argument: 72,
		},

		// Footer Max Line Rule
		(&rule.FooterMaxLineLenRule{}).Name(): {
			Argument: 72,
		},

		// Types Enum Rule
		(&rule.TypeEnumRule{}).Name(): {
			Argument: []interface{}{
				"feat", "fix", "docs", "style", "refactor", "perf",
				"test", "build", "ci", "chore", "revert",
			},
		},

		// Scope Enum Rule
		(&rule.ScopeEnumRule{}).Name(): {
			Argument: []interface{}{},
			Flags: map[string]interface{}{
				"allow-empty": true,
			},
		},

		// Body Min Len Rule
		(&rule.BodyMinLenRule{}).Name(): {
			Argument: 0,
		},

		// Body Max Len Rule
		(&rule.BodyMaxLenRule{}).Name(): {
			Argument: -1,
		},

		// Footer Min Len Rule
		(&rule.FooterMinLenRule{}).Name(): {
			Argument: 0,
		},

		// Footer Max Len Rule
		(&rule.FooterMaxLenRule{}).Name(): {
			Argument: -1,
		},

		// Type Min Len Rule
		(&rule.TypeMinLenRule{}).Name(): {
			Argument: 0,
		},

		// Type Max Len Rule
		(&rule.TypeMaxLenRule{}).Name(): {
			Argument: -1,
		},

		// Scope Min Len Rule
		(&rule.ScopeMinLenRule{}).Name(): {
			Argument: 0,
		},

		// Scope Max Len Rule
		(&rule.ScopeMaxLenRule{}).Name(): {
			Argument: -1,
		},

		// Description Min Len Rule
		(&rule.DescriptionMinLenRule{}).Name(): {
			Argument: 0,
		},

		// Description Max Len Rule
		(&rule.DescriptionMaxLenRule{}).Name(): {
			Argument: -1,
		},

		// Type Charset Rule
		(&rule.TypeCharsetRule{}).Name(): {
			Argument: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},

		// Scope Charset Rule
		(&rule.ScopeCharsetRule{}).Name(): {
			Argument: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ/,",
		},

		// Footer Enum Rule
		(&rule.FooterEnumRule{}).Name(): {
			Argument: []interface{}{},
		},
	}

	def := &lint.Config{
		MinVersion: internal.Version(),
		Formatter:  (&formatter.DefaultFormatter{}).Name(),
		Rules:      rules,
		Severity:   severity,
		Settings:   settings,
	}
	return def
}

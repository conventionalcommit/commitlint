package lint

// RuleSetting represent config for a rule
type RuleSetting struct {
	Argument interface{}            `yaml:"argument"`
	Flags    map[string]interface{} `yaml:"flags"`
}

// SeverityConfig represent severity levels for rules
type SeverityConfig struct {
	Default Severity            `yaml:"default"`
	Rules   map[string]Severity `yaml:"rules"`
}

// Config represent linter config
type Config struct {
	// MinVersion is the minimum version of commitlint required
	// should be in semver format
	MinVersion string `yaml:"min-version"`

	// DeprecatedVersion is the old "version" key, kept for backward compatibility.
	// Use MinVersion ("min-version") in new config files.
	DeprecatedVersion string `yaml:"version,omitempty"`

	// Formatter of the lint result
	Formatter string `yaml:"formatter"`

	// Enabled Rules
	Rules []string `yaml:"rules"`

	// Severity
	Severity SeverityConfig `yaml:"severity"`

	// Settings is rule name to rule settings
	Settings map[string]RuleSetting `yaml:"settings"`

	// DisableDefaultIgnores disables the built-in ignore patterns
	// (merge, revert, fixup, squash, etc.) when set to true.
	DisableDefaultIgnores bool `yaml:"disable-default-ignores"`

	// IgnorePatterns is a list of user-defined regex patterns.
	// If the first line of the commit message matches any pattern,
	// linting is skipped. These are added on top of the default
	// patterns (unless DisableDefaultIgnores is true).
	IgnorePatterns []string `yaml:"ignores"`

	// DefaultIgnorePatterns holds the built-in patterns (set by config package).
	// Not serialized to YAML - users never set this directly.
	DefaultIgnorePatterns []string `yaml:"-"`
}

// EffectiveIgnorePatterns returns the combined list of patterns the linter should use.
// If DisableDefaultIgnores is true, only user-defined patterns are returned.
func (c *Config) EffectiveIgnorePatterns() []string {
	if c.DisableDefaultIgnores {
		return c.IgnorePatterns
	}
	combined := make([]string, 0, len(c.DefaultIgnorePatterns)+len(c.IgnorePatterns))
	combined = append(combined, c.DefaultIgnorePatterns...)
	combined = append(combined, c.IgnorePatterns...)
	return combined
}

// GetRule returns RuleConfig for given rule name
func (c *Config) GetRule(ruleName string) RuleSetting {
	return c.Settings[ruleName]
}

// GetSeverity returns Severity for given ruleName
func (c *Config) GetSeverity(ruleName string) Severity {
	s, ok := c.Severity.Rules[ruleName]
	if ok {
		return s
	}
	return c.Severity.Default
}

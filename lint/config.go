package lint

// RuleSetting represent config for a rule
type RuleSetting struct {
	Argument interface{}            `yaml:"argument"`
	Flags    map[string]interface{} `yaml:"flags,omitempty"`
}

// SeverityConfig represent severity levels for rules
type SeverityConfig struct {
	Default Severity            `yaml:"default"`
	Rules   map[string]Severity `yaml:"rules,omitempty"`
}

// Config represent linter config
type Config struct {
	// MinVersion is the minimum version of commitlint required
	// should be in semver format
	MinVersion string `yaml:"version"`

	// Formatter of the lint result
	Formatter string `yaml:"formatter"`

	// Enabled Rules
	Rules []string `yaml:"rules"`

	// Severity
	Severity SeverityConfig `yaml:"severity"`

	// Settings is rule name to rule settings
	Settings map[string]RuleSetting `yaml:"settings"`
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

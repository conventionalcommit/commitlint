package lint

// RuleConfig represent config for a rule
type RuleConfig struct {
	Enabled bool `yaml:"enabled"`

	Severity Severity `yaml:"severity"`

	Argument interface{} `yaml:"argument"`

	// Flags are optional key value pairs
	Flags map[string]interface{} `yaml:"flags"`
}

// Config represent linter config
type Config struct {
	// MinVersion is the minimum version of commitlint required
	// should be in semver format
	MinVersion string `yaml:"version"`

	// Formatter of the lint result
	Formatter string `yaml:"formatter"`

	// Rules is rule name to rule config map
	Rules map[string]RuleConfig `yaml:"rules"`
}

// GetRule returns RuleConfig for given rule name
func (c *Config) GetRule(ruleName string) RuleConfig {
	return c.Rules[ruleName]
}

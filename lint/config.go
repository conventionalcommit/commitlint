package lint

// Rule Severity Constants
const (
	SeverityWarn  = "warn"
	SeverityError = "error"
)

// Config represent linter config
type Config struct {
	Formatter string                `yaml:"formatter"`
	Rules     map[string]RuleConfig `yaml:"rules"`
}

// RuleConfig represent config for a rule
type RuleConfig struct {
	Enabled  bool        `yaml:"enabled"`
	Severity string      `yaml:"severity"`
	Argument interface{} `yaml:"argument"`

	// Optional flags
	Flags map[string]interface{} `yaml:"flags"`
}

// GetRule returns RuleConfig for given ruleName
func (c *Config) GetRule(ruleName string) RuleConfig {
	return c.Rules[ruleName]
}

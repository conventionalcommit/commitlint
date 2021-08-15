package commitlint

// Rule Type Constants
const (
	SeverityWarn  = "warn"
	SeverityError = "error"
)

// Config represent the rules for commit message
type Config struct {
	Header Header `yaml:"header"`
	Body   Body   `yaml:"body"`
	Footer Footer `yaml:"footer"`
}

// Header contains rules for Header in commit message
type Header struct {
	MinLength IntConf  `yaml:"min-length"`
	MaxLength IntConf  `yaml:"max-length"`
	Types     EnumConf `yaml:"types"`
	Scopes    EnumConf `yaml:"scopes"`
}

// Body contains rules for Body in commit message
type Body struct {
	CanBeEmpty    bool    `yaml:"can-be-empty"`
	MaxLineLength IntConf `yaml:"max-line-length"`
}

// Footer contains rules for Footer in commit message
type Footer struct {
	CanBeEmpty    bool    `yaml:"can-be-empty"`
	MaxLineLength IntConf `yaml:"max-line-length"`
}

// IntConf represent int config
type IntConf struct {
	Enabled bool   `yaml:"enabled"`
	Type    string `yaml:"type"`
	Value   int    `yaml:"value"`
}

// EnumConf represent enums config
type EnumConf struct {
	Enabled bool     `yaml:"enabled"`
	Type    string   `yaml:"type"`
	Value   []string `yaml:"value"`
}

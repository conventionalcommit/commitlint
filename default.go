package commitlint

var defConf = &Config{
	Header: Header{
		MinLength: IntConf{Enabled: true, Type: SeverityError, Value: 10},
		MaxLength: IntConf{Enabled: true, Type: SeverityError, Value: 50},
		Scopes: EnumConf{
			Enabled: false,
			Type:    SeverityError,
			Value:   []string{},
		},
		Types: EnumConf{
			Enabled: true,
			Type:    SeverityError,
			Value: []string{
				"feat",
				"fix",
				"docs",
				"style",
				"refactor",
				"perf",
				"test",
				"build",
				"ci",
				"chore",
				"revert",
				"merge",
			},
		},
	},
	Body: Body{
		CanBeEmpty:    true,
		MaxLineLength: IntConf{Enabled: true, Type: SeverityError, Value: 72},
	},
	Footer: Footer{
		CanBeEmpty:    true,
		MaxLineLength: IntConf{Enabled: true, Type: SeverityError, Value: 72},
	},
}

// NewDefaultLinter returns Linter with default config
func NewDefaultLinter() *Linter {
	return &Linter{conf: defConf}
}

// DefaultConfToFile writes default config to given file
func DefaultConfToFile(confPath string) error {
	return WriteConfToFile(confPath, defConf)
}

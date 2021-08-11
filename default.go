package commitlint

import (
	"bufio"
	"os"

	"gopkg.in/yaml.v3"
)

var defConf = &Config{
	Header: Header{
		MinLength: IntConf{Enabled: true, Type: ErrorType, Value: 10},
		MaxLength: IntConf{Enabled: true, Type: ErrorType, Value: 50},
		Scopes: EnumConf{
			Enabled: false,
			Type:    ErrorType,
			Value:   []string{},
		},
		Types: EnumConf{
			Enabled: true,
			Type:    ErrorType,
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
		MaxLineLength: IntConf{Enabled: true, Type: ErrorType, Value: 72},
	},
	Footer: Footer{
		CanBeEmpty:    true,
		MaxLineLength: IntConf{Enabled: true, Type: ErrorType, Value: 72},
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

// WriteConfToFile util func to write config object to given file
func WriteConfToFile(confPath string, conf *Config) (retErr error) {
	file, err := os.Create(confPath)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	w := bufio.NewWriter(file)
	defer func() {
		err1 := w.Flush()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	enc := yaml.NewEncoder(w)
	return enc.Encode(conf)
}

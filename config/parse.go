package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/conventionalcommit/commitlint/lint"
)

// Parse parse Config from given file
func Parse(confPath string) (*lint.Config, error) {
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	conf := &lint.Config{}
	err = yaml.Unmarshal(confBytes, conf)
	if err != nil {
		return nil, err
	}

	err = Validate(conf)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return conf, nil
}

// Validate parses Config from given data
func Validate(conf *lint.Config) error {
	if conf.Formatter == "" {
		return errors.New("formatter is empty")
	}

	for ruleName, r := range conf.Rules {
		switch r.Severity {
		case lint.SeverityError:
		case lint.SeverityWarn:
		default:
			return fmt.Errorf("unknown severity level '%s' for rule '%s'", r.Severity, ruleName)
		}
	}

	return nil
}

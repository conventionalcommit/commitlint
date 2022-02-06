// Package config contains helpers, defaults for linter
package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/mod/semver"
	yaml "gopkg.in/yaml.v2"

	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/internal/registry"
	"github.com/conventionalcommit/commitlint/lint"
)

// Parse parse given file in confPath, and return Config instance, error if any
func Parse(confPath string) (*lint.Config, error) {
	confPath = filepath.Clean(confPath)
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	conf := &lint.Config{}
	err = yaml.UnmarshalStrict(confBytes, conf)
	if err != nil {
		return nil, err
	}

	if conf.Formatter == "" {
		return nil, errors.New("conf error: formatter is empty")
	}

	err = isValidVersion(conf.MinVersion)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Validate validates given config instance, it checks the following
// If formatters, rules are registered/known
// If arguments to rules are valid
// If version is valid and atleast minimum than commitlint version used
func Validate(conf *lint.Config) []error {
	var errs []error

	err := isValidVersion(conf.MinVersion)
	if err != nil {
		errs = append(errs, err)
	}

	if conf.Formatter == "" {
		errs = append(errs, errors.New("formatter is empty"))
	} else {
		_, ok := registry.GetFormatter(conf.Formatter)
		if !ok {
			errs = append(errs, fmt.Errorf("unknown formatter '%s'", conf.Formatter))
		}
	}

	for ruleName, r := range conf.Rules {
		// Check Severity Level of rule config
		switch r.Severity {
		case lint.SeverityError, lint.SeverityWarn:
		default:
			errs = append(errs, fmt.Errorf("unknown severity level '%s' for rule '%s'", r.Severity, ruleName))
		}

		// Check if rule is registered
		ruleData, ok := registry.GetRule(ruleName)
		if !ok {
			errs = append(errs, fmt.Errorf("unknown rule '%s'", ruleName))
			continue
		}

		err := ruleData.Apply(r.Argument, r.Flags)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// LookupAndParse gets the config path according to the precedence
// if exists, parses the config file and returns config instance
func LookupAndParse() (*lint.Config, error) {
	confFilePath, confType, err := internal.LookupConfigPath()
	if err != nil {
		return nil, err
	}

	if confType == internal.DefaultConfig {
		return Default(), nil
	}

	conf, err := Parse(confFilePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// WriteTo writes config in yaml format to given io.Writer
func WriteTo(w io.Writer, conf *lint.Config) (retErr error) {
	enc := yaml.NewEncoder(w)
	defer func() {
		err := enc.Close()
		if retErr == nil && err != nil {
			retErr = err
		}
	}()
	return enc.Encode(conf)
}

func isValidVersion(versionNo string) error {
	if versionNo == "" {
		return errors.New("version is empty")
	}
	if !semver.IsValid(versionNo) {
		return errors.New("invalid version should be in semver format")
	}
	return nil
}

func checkIfMinVersion(versionNo string) error {
	cmp := semver.Compare(internal.Version(), versionNo)
	if cmp != -1 {
		return nil
	}
	return fmt.Errorf("min version required is %s. you have %s.\nupgrade commitlint", versionNo, internal.Version())
}

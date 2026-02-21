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

	"github.com/conventionalcommit/commitlint/formatter"
	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/internal/registry"
	"github.com/conventionalcommit/commitlint/lint"
)

// Parse parse given file in confPath, and return Config instance, error if any
func Parse(confPath string) (*lint.Config, error) {
	confPath = filepath.Clean(confPath)
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("config file error: %w", err)
	}

	conf := &lint.Config{
		Formatter: (&formatter.DefaultFormatter{}).Name(),
		Severity: lint.SeverityConfig{
			Default: lint.SeverityError,
		},
	}

	err = yaml.UnmarshalStrict(confBytes, conf)
	if err != nil {
		return nil, fmt.Errorf("config file error: %w", err)
	}

	// Backward compatibility: accept old "version" key
	if conf.MinVersion == "" && conf.DeprecatedVersion != "" {
		conf.MinVersion = conf.DeprecatedVersion
	}
	conf.DeprecatedVersion = ""

	// Default to current version if neither key was provided
	if conf.MinVersion == "" {
		conf.MinVersion = internal.Version()
	}

	if conf.Formatter == "" {
		return nil, errors.New("config error: formatter is empty")
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
// If version is valid and at least minimum than commitlint version used
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

	// Check Severity Level
	if !isSeverityValid(conf.Severity.Default) {
		errs = append(errs, fmt.Errorf("unknown default severity level '%s'", conf.Severity.Default))
	}

	for ruleName, sev := range conf.Severity.Rules {
		// Check Severity Level of rule config
		if !isSeverityValid(sev) {
			errs = append(errs, fmt.Errorf("unknown severity level '%s' for rule '%s'", sev, ruleName))
		}
	}

	for _, ruleName := range conf.Rules {
		// Check if rule is registered
		_, ok := registry.GetRule(ruleName)
		if !ok {
			errs = append(errs, fmt.Errorf("unknown rule '%s'", ruleName))
			continue
		}
	}

	// Check for duplicate rules
	ruleSeen := make(map[string]struct{}, len(conf.Rules))
	for _, ruleName := range conf.Rules {
		if _, exists := ruleSeen[ruleName]; exists {
			errs = append(errs, fmt.Errorf("duplicate rule '%s' in rules list", ruleName))
		} else {
			ruleSeen[ruleName] = struct{}{}
		}
	}

	for ruleName, ruleSetting := range conf.Settings {
		// Check if rule is registered
		ruleData, ok := registry.GetRule(ruleName)
		if !ok {
			errs = append(errs, fmt.Errorf("unknown rule '%s'", ruleName))
			continue
		}

		err := ruleData.Apply(ruleSetting)
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
		return NewDefault(), nil
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

func isSeverityValid(s lint.Severity) bool {
	return s == lint.SeverityError || s == lint.SeverityWarn
}

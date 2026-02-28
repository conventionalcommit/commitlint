package config

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/registry"
	"golang.org/x/mod/semver"
)

// NewLinter returns Linter for given confFilePath
func NewLinter(conf *lint.Config) (*lint.Linter, error) {
	err := checkIfMinVersion(conf.MinVersion)
	if err != nil {
		return nil, err
	}

	rules, err := GetEnabledRules(conf)
	if err != nil {
		return nil, err
	}

	return lint.New(conf, rules)
}

// GetFormatter returns the formatter as defined in conf
func GetFormatter(conf *lint.Config) (lint.Formatter, error) {
	err := checkIfMinVersion(conf.MinVersion)
	if err != nil {
		return nil, err
	}

	format, ok := registry.GetFormatter(conf.Formatter)
	if !ok {
		return nil, fmt.Errorf("config error: '%s' formatter not found", conf.Formatter)
	}
	return format, nil
}

// GetEnabledRules forms Rule object for rules which are enabled in config
func GetEnabledRules(conf *lint.Config) ([]lint.Rule, error) {
	enabledRules := make([]lint.Rule, 0, len(conf.Rules))

	// To check if duplicate rule is added
	addedRules := make(map[string]struct{})

	for _, ruleName := range conf.Rules {
		if _, ok := addedRules[ruleName]; ok {
			continue
		}

		// Checking if rule is registered
		// before checking if rule is enabled
		r, ok := registry.GetRule(ruleName)
		if !ok {
			return nil, fmt.Errorf("config error: '%s' rule not found", ruleName)
		}

		rConf, ok := conf.Settings[ruleName]
		if !ok {
			return nil, fmt.Errorf("config error: '%s' rule settings not found", ruleName)
		}

		err := r.Apply(rConf)
		if err != nil {
			return nil, fmt.Errorf("config error: %v", err)
		}
		enabledRules = append(enabledRules, r)
		addedRules[r.Name()] = struct{}{}
	}

	return enabledRules, nil
}

// ValidateLint validates given lint config instance, it checks the following
// If formatters, rules are registered/known
// If arguments to rules are valid
// If version is valid and at least minimum than commitlint version used
func ValidateLint(conf *lint.Config) []error {
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

	// Validate ignore patterns (both default and user-defined)
	for _, pattern := range conf.EffectiveIgnorePatterns() {
		_, err := regexp.Compile(pattern)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid ignore pattern %q: %w", pattern, err))
		}
	}

	return errs
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

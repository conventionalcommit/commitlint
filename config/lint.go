package config

import (
	"fmt"

	"github.com/conventionalcommit/commitlint/internal/registry"
	"github.com/conventionalcommit/commitlint/lint"
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

	for _, ruleName := range conf.Rules {
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
	}

	return enabledRules, nil
}

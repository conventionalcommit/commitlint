package config

import (
	"fmt"

	"github.com/conventionalcommit/commitlint/lint"
)

// GetLinter returns Linter for given confFilePath
func GetLinter(conf *lint.Config) (*lint.Linter, error) {
	rules, err := GetEnabledRules(conf)
	if err != nil {
		return nil, err
	}
	return lint.New(conf, rules)
}

// GetFormatter returns the formatter as defined in conf
func GetFormatter(conf *lint.Config) (lint.Formatter, error) {
	format, ok := globalRegistry.GetFormatter(conf.Formatter)
	if !ok {
		return nil, fmt.Errorf("'%s' formatter not found", conf.Formatter)
	}
	return format, nil
}

// GetEnabledRules forms Rule object for rules which are enabled in config
func GetEnabledRules(conf *lint.Config) ([]lint.Rule, error) {
	enabledRules := make([]lint.Rule, 0, len(conf.Rules))

	for ruleName, ruleConfig := range conf.Rules {
		// Checking if rule is registered
		// before checking if rule is enabled
		r, ok := globalRegistry.GetRule(ruleName)
		if !ok {
			return nil, fmt.Errorf("'%s' rule not found", ruleName)
		}

		if !ruleConfig.Enabled {
			continue
		}

		err := r.Apply(ruleConfig.Argument, ruleConfig.Flags)
		if err != nil {
			return nil, err
		}
		enabledRules = append(enabledRules, r)
	}

	return enabledRules, nil
}

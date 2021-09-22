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
	return lint.NewLinter(conf, rules)
}

// GetFormatter returns the formatter as defined in conf
func GetFormatter(c *lint.Config) (lint.Formatter, error) {
	for _, f := range defaultFormatters {
		if f.Name() == c.Formatter {
			return f, nil
		}
	}
	return nil, fmt.Errorf("%s formatter not found", c.Formatter)
}

// GetEnabledRules forms Rule object for rules which are enabled in config
func GetEnabledRules(conf *lint.Config) ([]lint.Rule, error) {
	// rules lookup map
	rulesMap := map[string]lint.Rule{}
	for _, r := range defaultRules {
		rulesMap[r.Name()] = r
	}

	enabledRules := make([]lint.Rule, 0, len(conf.Rules))

	for ruleName, ruleConfig := range conf.Rules {
		r, ok := rulesMap[ruleName]
		if !ok {
			return nil, fmt.Errorf("unknown rule: %s", ruleName)
		}
		if ruleConfig.Enabled {
			err := r.Apply(ruleConfig.Argument, ruleConfig.Flags)
			if err != nil {
				return nil, err
			}
			enabledRules = append(enabledRules, r)
		}
	}

	return enabledRules, nil
}

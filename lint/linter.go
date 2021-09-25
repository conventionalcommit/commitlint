// Package lint provides a simple linter for conventional commits
package lint

import (
	"github.com/conventionalcommit/commitlint/message"
)

// Linter is linter for commit message
type Linter struct {
	conf  *Config
	rules []Rule
}

// NewLinter returns a new Linter instance with given config and rules
func NewLinter(conf *Config, rules []Rule) (*Linter, error) {
	return &Linter{conf: conf, rules: rules}, nil
}

// Lint checks the given commitMsg string against rules
func (l *Linter) Lint(commitMsg string) (*Result, error) {
	msg, isHeaderErr, err := message.Parse(commitMsg)
	if err != nil {
		if isHeaderErr {
			return l.headerErrorRule(commitMsg), nil
		}
		return nil, err
	}
	return l.LintCommit(msg)
}

// LintCommit checks the given message.Commit against rules
func (l *Linter) LintCommit(msg *message.Commit) (*Result, error) {
	res := newResult(msg.FullCommit)

	for _, rule := range l.rules {
		result, isOK := rule.Validate(msg)
		if !isOK {
			ruleConf := l.conf.GetRule(rule.Name())
			res.add(RuleResult{
				Name:     rule.Name(),
				Severity: ruleConf.Severity,
				Message:  result,
			})
		}
	}

	return res, nil
}

func (l *Linter) headerErrorRule(commitMsg string) *Result {
	// TODO: show more information
	res := newResult(commitMsg)
	res.add(RuleResult{
		Name:     "parser",
		Severity: SeverityError,
		Message:  "commit header is not valid",
	})
	return res
}

// Package lint provides a simple linter for conventional commits
package lint

import (
	"github.com/conventionalcommit/parser"
)

// Linter is linter for commit message
type Linter struct {
	conf  *Config
	rules []Rule

	parser Parser
}

// New returns a new Linter instance with given config and rules
func New(conf *Config, rules []Rule) (*Linter, error) {
	l := &Linter{
		conf:   conf,
		rules:  rules,
		parser: parser.New(),
	}
	return l, nil
}

// Lint checks the given commitMsg string against rules
func (l *Linter) Lint(commitMsg string) (*Failure, error) {
	msg, err := l.parser.Parse(commitMsg)
	if err != nil {
		return l.parserErrorRule(commitMsg, err)
	}
	return l.LintCommit(msg)
}

// LintCommit checks the given Commit against rules
func (l *Linter) LintCommit(msg *Commit) (*Failure, error) {
	res := newFailure(msg.Message())

	for _, rule := range l.rules {
		currentRule := rule
		ruleConf := l.conf.GetRule(currentRule.Name())
		ruleRes, isValid := l.runRule(currentRule, ruleConf.Severity, msg)
		if !isValid {
			res.add(ruleRes)
		}
	}

	return res, nil
}

func (l *Linter) runRule(rule Rule, severity Severity, msg *Commit) (*RuleFailure, bool) {
	failMsgs, isOK := rule.Validate(msg)
	if isOK {
		return nil, true
	}
	res := newRuleFailure(rule.Name(), failMsgs, severity)
	return res, false
}

func (l *Linter) parserErrorRule(commitMsg string, err error) (*Failure, error) {
	res := newFailure(commitMsg)

	errMsg := err.Error()

	ruleFail := newRuleFailure("parser", []string{errMsg}, SeverityError)
	res.add(ruleFail)

	return res, nil
}

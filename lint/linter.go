// Package lint provides a simple linter for conventional commits
package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Linter is linter for commit message
type Linter struct {
	conf  *Config
	rules []Rule

	parser         Parser
	ignorePatterns []*regexp.Regexp
}

// New returns a new Linter instance with given config and rules
func New(conf *Config, rules []Rule) (*Linter, error) {
	compiled, err := compilePatterns(conf.EffectiveIgnorePatterns())
	if err != nil {
		return nil, err
	}

	l := &Linter{
		conf:           conf,
		rules:          rules,
		parser:         newParser(),
		ignorePatterns: compiled,
	}
	return l, nil
}

// ParseAndLint checks the given commitMsg string against rules
func (l *Linter) ParseAndLint(commitMsg string) (*Result, error) {
	if l.isIgnored(commitMsg) {
		return newResult(commitMsg), nil
	}

	msg, err := l.parser.Parse(commitMsg)
	if err != nil {
		issues := l.parserErrorRule(commitMsg, err)
		return newResult(commitMsg, issues...), nil
	}
	return l.Lint(msg)
}

// isIgnored checks if the first line of the commit message
// matches any of the configured ignore patterns
func (l *Linter) isIgnored(commitMsg string) bool {
	if len(l.ignorePatterns) == 0 {
		return false
	}

	firstLine := strings.Split(commitMsg, "\n")[0]
	for _, re := range l.ignorePatterns {
		if re.MatchString(firstLine) {
			return true
		}
	}
	return false
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid ignore pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return compiled, nil
}

// Lint checks the given Commit against rules
func (l *Linter) Lint(msg Commit) (*Result, error) {
	issues := make([]*Issue, 0, len(l.rules))

	for _, rule := range l.rules {
		currentRule := rule
		severity := l.conf.GetSeverity(currentRule.Name())
		issue, isValid := l.runRule(currentRule, severity, msg)
		if !isValid {
			issues = append(issues, issue)
		}
	}

	return newResult(msg.Message(), issues...), nil
}

func (l *Linter) runRule(rule Rule, severity Severity, msg Commit) (*Issue, bool) {
	issue, isValid := rule.Validate(msg)
	if isValid {
		return nil, true
	}

	issue.ruleName = rule.Name()
	issue.severity = severity
	return issue, false
}

func (l *Linter) parserErrorRule(commitMsg string, err error) []*Issue {
	issue := NewIssue(err.Error())
	issue.ruleName = "parser"
	issue.severity = SeverityError
	return []*Issue{issue}
}

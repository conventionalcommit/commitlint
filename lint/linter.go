// Package lint provides a simple linter for conventional commits
package lint

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/conventionalcommit/commitlint/commit"
)

// Linter is linter for commit message
type Linter struct {
	conf   *Config
	rules  []Rule
	parser commit.Parser

	ignorePatterns []*regexp.Regexp
}

// New returns a new Linter instance with given config and rules using the
// default inspection parser.
func New(conf *Config, rules []Rule) (*Linter, error) {
	return newLinter(conf, rules, nil)
}

// NewWithParser returns a new Linter instance using parser for ParseAndLint.
func NewWithParser(conf *Config, rules []Rule, parser commit.Parser) (*Linter, error) {
	if parser == nil {
		return nil, errors.New("parser cannot be nil")
	}
	return newLinter(conf, rules, parser)
}

func newLinter(conf *Config, rules []Rule, parser commit.Parser) (*Linter, error) {
	compiled, err := compilePatterns(conf.EffectiveIgnorePatterns())
	if err != nil {
		return nil, err
	}

	return &Linter{
		conf:           conf,
		rules:          rules,
		parser:         parser,
		ignorePatterns: compiled,
	}, nil
}

// ParseAndLint checks the given commitMsg string against rules
func (l *Linter) ParseAndLint(commitMsg string) (*Result, error) {
	if l.isIgnored(commitMsg) {
		return newResult(commitMsg), nil
	}

	var (
		msg         Commit
		diagnostics []error
	)
	if l.parser != nil {
		var err error
		msg, err = l.parser.Parse(commitMsg)
		if err != nil {
			diagnostics = []error{err}
		}
	} else {
		msg, diagnostics = commit.Inspect(commitMsg)
	}

	if msg != nil && len(diagnostics) == 0 {
		return l.Lint(msg)
	}

	var issues []*Issue
	for _, diagnostic := range diagnostics {
		issues = append(issues, l.parserErrorRule(commitMsg, diagnostic)...)
	}
	if msg == nil {
		return newResult(commitMsg, issues...), nil
	}

	result, err := l.Lint(msg)
	if err != nil {
		return nil, err
	}
	result.issues = append(issues, result.issues...)
	return result, nil
}

// isIgnored checks if the first line of the commit message
// matches any of the configured ignore patterns
func (l *Linter) isIgnored(commitMsg string) bool {
	if len(l.ignorePatterns) == 0 {
		return false
	}

	firstLine := commit.FirstLine(commitMsg)
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
	var infos []string
	if line, col, ok := commit.DiagnosticPosition(commitMsg, err); ok {
		infos = []string{fmt.Sprintf("line %d, column %d", line, col)}
	}
	issue := NewIssue(err.Error(), infos...)
	issue.ruleName = "parser"
	issue.severity = SeverityError
	return []*Issue{issue}
}

// Package lint provides a simple linter for conventional commits
package lint

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
		parser: newParser(),
	}
	return l, nil
}

// ParseAndLint checks the given commitMsg string against rules
func (l *Linter) ParseAndLint(commitMsg string) (*Result, error) {
	msg, err := l.parser.Parse(commitMsg)
	if err != nil {
		issues := l.parserErrorRule(commitMsg, err)
		return newResult(commitMsg, issues...), nil
	}
	return l.Lint(msg)
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

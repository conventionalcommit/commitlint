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

// Lint checks the given commitMsg string against rules
func (l *Linter) Lint(commitMsg string) (*Result, error) {
	msg, err := l.parser.Parse(commitMsg)
	if err != nil {
		return l.parserErrorRule(commitMsg, err)
	}
	return l.LintCommit(msg)
}

// LintCommit checks the given Commit against rules
func (l *Linter) LintCommit(msg Commit) (*Result, error) {
	res := newResult(msg.Message())

	for _, rule := range l.rules {
		currentRule := rule
		severity := l.conf.GetSeverity(currentRule.Name())
		ruleRes, isValid := l.runRule(currentRule, severity, msg)
		if !isValid {
			res.add(ruleRes)
		}
	}

	return res, nil
}

func (l *Linter) runRule(rule Rule, severity Severity, msg Commit) (*Issue, bool) {
	failMsgs, isOK := rule.Validate(msg)
	if isOK {
		return nil, true
	}
	res := newIssue(rule.Name(), failMsgs, severity)
	return res, false
}

func (l *Linter) parserErrorRule(commitMsg string, err error) (*Result, error) {
	res := newResult(commitMsg)

	errMsg := err.Error()

	ruleFail := newIssue("parser", []string{errMsg}, SeverityError)
	res.add(ruleFail)

	return res, nil
}

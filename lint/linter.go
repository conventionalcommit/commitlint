// Package lint provides a simple linter for conventional commits
package lint

// Linter is linter for commit message
type Linter struct {
	conf  *Config
	rules []Rule
}

// New returns a new Linter instance with given config and rules
func New(conf *Config, rules []Rule) (*Linter, error) {
	return &Linter{conf: conf, rules: rules}, nil
}

// Lint checks the given commitMsg string against rules
func (l *Linter) Lint(commitMsg string) (*Result, error) {
	msg, err := Parse(commitMsg)
	if err != nil {
		if isHeaderErr(err) {
			return l.headerErrorRule(commitMsg), nil
		}
		return nil, err
	}
	return l.LintCommit(msg)
}

// LintCommit checks the given Commit against rules
func (l *Linter) LintCommit(msg *Commit) (*Result, error) {
	res := newResult(msg.FullCommit)

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

func (l *Linter) runRule(rule Rule, severity Severity, msg *Commit) (RuleResult, bool) {
	ruleMsg, isOK := rule.Validate(msg)
	if isOK {
		return RuleResult{}, true
	}
	res := RuleResult{
		Name:     rule.Name(),
		Severity: severity,
		Message:  ruleMsg,
	}
	return res, false
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

package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*FooterMaxLineLenRule)(nil)

// FooterMaxLineLenRule to validate max line length of footer
type FooterMaxLineLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *FooterMaxLineLenRule) Name() string { return "footer-max-line-length" }

// Apply sets the needed argument for the rule
func (r *FooterMaxLineLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates FooterMaxLineLenRule rule
func (r *FooterMaxLineLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLineLength("footer", r.CheckLen, msg.Footer())
}

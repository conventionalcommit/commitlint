package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*FooterMaxLenRule)(nil)

// FooterMaxLenRule to validate max length of footer
type FooterMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *FooterMaxLenRule) Name() string { return "footer-max-length" }

// Apply sets the needed argument for the rule
func (r *FooterMaxLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates FooterMaxLenRule
func (r *FooterMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("footer", r.CheckLen, msg.Footer())
}

package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*HeadMaxLenRule)(nil)

// HeadMaxLenRule to validate max length of header
type HeadMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *HeadMaxLenRule) Name() string { return "header-max-length" }

// Apply sets the needed argument for the rule
func (r *HeadMaxLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates HeadMaxLenRule
func (r *HeadMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("header", r.CheckLen, msg.Header())
}

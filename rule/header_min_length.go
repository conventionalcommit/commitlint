package rule

import (
	"github.com/conventionalcommit/commitlint/lint"
)

var _ lint.Rule = (*HeadMinLenRule)(nil)

// HeadMinLenRule to validate min length of header
type HeadMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *HeadMinLenRule) Name() string { return "header-min-length" }

// Apply sets the needed argument for the rule
func (r *HeadMinLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates HeadMinLenRule
func (r *HeadMinLenRule) Validate(msg lint.Commit) ([]string, bool) {
	return checkMinLen(r.CheckLen, msg.Header())
}

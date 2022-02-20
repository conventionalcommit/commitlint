package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*TypeMinLenRule)(nil)

// TypeMinLenRule to validate min length of type
type TypeMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *TypeMinLenRule) Name() string { return "type-min-length" }

// Apply sets the needed argument for the rule
func (r *TypeMinLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates TypeMinLenRule
func (r *TypeMinLenRule) Validate(msg lint.Commit) ([]string, bool) {
	return checkMinLen(r.CheckLen, msg.Type())
}

package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*TypeMaxLenRule)(nil)

// TypeMaxLenRule to validate max length of type
type TypeMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *TypeMaxLenRule) Name() string { return "type-max-length" }

// Apply sets the needed argument for the rule
func (r *TypeMaxLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates TypeMaxLenRule
func (r *TypeMaxLenRule) Validate(msg lint.Commit) ([]string, bool) {
	return checkMaxLen(r.CheckLen, msg.Type())
}

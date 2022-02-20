package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*BodyMaxLenRule)(nil)

// BodyMaxLenRule to validate max length of body
type BodyMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *BodyMaxLenRule) Name() string { return "body-max-length" }

// Apply sets the needed argument for the rule
func (r *BodyMaxLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates BodyMaxLenRule
func (r *BodyMaxLenRule) Validate(msg lint.Commit) ([]string, bool) {
	return checkMaxLen(r.CheckLen, msg.Body())
}

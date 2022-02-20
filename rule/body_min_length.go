package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*BodyMinLenRule)(nil)

// BodyMinLenRule to validate min length of body
type BodyMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *BodyMinLenRule) Name() string { return "body-min-length" }

// Apply sets the needed argument for the rule
func (r *BodyMinLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates BodyMinLenRule
func (r *BodyMinLenRule) Validate(msg lint.Commit) ([]string, bool) {
	return checkMinLen(r.CheckLen, msg.Body())
}

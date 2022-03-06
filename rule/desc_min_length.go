package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*DescriptionMinLenRule)(nil)

// DescriptionMinLenRule to validate min length of description
type DescriptionMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *DescriptionMinLenRule) Name() string { return "description-min-length" }

// Apply sets the needed argument for the rule
func (r *DescriptionMinLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates DescriptionMinLenRule
func (r *DescriptionMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("description", r.CheckLen, msg.Description())
}

package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*DescriptionMaxLenRule)(nil)

// DescriptionMaxLenRule to validate max length of type
type DescriptionMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *DescriptionMaxLenRule) Name() string { return "description-max-length" }

// Apply sets the needed argument for the rule
func (r *DescriptionMaxLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates DescriptionMaxLenRule
func (r *DescriptionMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("description", r.CheckLen, msg.Description())
}

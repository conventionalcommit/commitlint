package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*ScopeMaxLenRule)(nil)

// ScopeMaxLenRule to validate max length of type
type ScopeMaxLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *ScopeMaxLenRule) Name() string { return "scope-max-length" }

// Apply sets the needed argument for the rule
func (r *ScopeMaxLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates ScopeMaxLenRule
func (r *ScopeMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("scope", r.CheckLen, msg.Scope())
}

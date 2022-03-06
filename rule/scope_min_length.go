package rule

import "github.com/conventionalcommit/commitlint/lint"

var _ lint.Rule = (*ScopeMinLenRule)(nil)

// ScopeMinLenRule to validate min length of scope
type ScopeMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *ScopeMinLenRule) Name() string { return "scope-min-length" }

// Apply sets the needed argument for the rule
func (r *ScopeMinLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates ScopeMinLenRule
func (r *ScopeMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("scope", r.CheckLen, msg.Scope())
}

package rule

import (
	"github.com/conventionalcommit/commitlint/lint"
)

var _ lint.Rule = (*ScopeCharsetRule)(nil)

// ScopeCharsetRule to validate max length of header
type ScopeCharsetRule struct {
	Charset string
}

// Name return name of the rule
func (r *ScopeCharsetRule) Name() string { return "scope-charset" }

// Apply sets the needed argument for the rule
func (r *ScopeCharsetRule) Apply(setting lint.RuleSetting) error {
	err := setStringArg(&r.Charset, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates ScopeCharsetRule
func (r *ScopeCharsetRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	invalidChars, isValid := validateCharset(r.Charset, msg.Scope())
	if isValid {
		return nil, true
	}

	desc := "type can only have these chars [" + r.Charset + "]"
	err := "invalid characters [" + invalidChars + "]"
	return lint.NewIssue(desc, err), false
}

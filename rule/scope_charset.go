package rule

import (
	"fmt"

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
func (r *ScopeCharsetRule) Validate(msg lint.Commit) ([]string, bool) {
	invalidChars, isValid := checkCharset(r.Charset, msg.Scope())
	if isValid {
		return nil, true
	}

	errMsg := fmt.Sprintf("scope contains invalid char '%s', allowed chars are [%s]", invalidChars, r.Charset)
	return []string{errMsg}, false
}

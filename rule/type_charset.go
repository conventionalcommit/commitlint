package rule

import (
	"github.com/conventionalcommit/commitlint/lint"
)

var _ lint.Rule = (*TypeCharsetRule)(nil)

// TypeCharsetRule to validate max length of header
type TypeCharsetRule struct {
	Charset string
}

// Name return name of the rule
func (r *TypeCharsetRule) Name() string { return "type-charset" }

// Apply sets the needed argument for the rule
func (r *TypeCharsetRule) Apply(setting lint.RuleSetting) error {
	err := setStringArg(&r.Charset, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates TypeCharsetRule
func (r *TypeCharsetRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	invalidChars, isValid := validateCharset(r.Charset, msg.Type())
	if isValid {
		return nil, true
	}

	desc := "type can only have chars [" + r.Charset + "]"
	info := "invalid characters [" + invalidChars + "]"

	return lint.NewIssue(desc, info), false
}

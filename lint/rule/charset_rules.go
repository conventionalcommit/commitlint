package rule

import "github.com/conventionalcommit/commitlint/lint"

// Compile-time interface checks
var (
	_ lint.Rule = (*TypeCharsetRule)(nil)
	_ lint.Rule = (*ScopeCharsetRule)(nil)
)

// applyStringArg is a shared helper that extracts a string argument from a RuleSetting.
func applyStringArg(dst *string, ruleName string, setting lint.RuleSetting) error {
	if err := setStringArg(dst, setting.Argument); err != nil {
		return errInvalidArg(ruleName, err)
	}
	return nil
}

// TypeCharsetRule to validate charset of type
type TypeCharsetRule struct{ Charset string }

func (r *TypeCharsetRule) Name() string { return "type-charset" }
func (r *TypeCharsetRule) Apply(s lint.RuleSetting) error {
	return applyStringArg(&r.Charset, r.Name(), s)
}

func (r *TypeCharsetRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	invalidChars, isValid := validateCharset(r.Charset, msg.Type())
	if isValid {
		return nil, true
	}
	return lint.NewIssue(
		"type can only have chars ["+r.Charset+"]",
		"invalid characters ["+invalidChars+"]",
	), false
}

// ScopeCharsetRule to validate charset of scope
type ScopeCharsetRule struct{ Charset string }

func (r *ScopeCharsetRule) Name() string { return "scope-charset" }
func (r *ScopeCharsetRule) Apply(s lint.RuleSetting) error {
	return applyStringArg(&r.Charset, r.Name(), s)
}

func (r *ScopeCharsetRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	invalidChars, isValid := validateCharset(r.Charset, msg.Scope())
	if isValid {
		return nil, true
	}
	return lint.NewIssue(
		"scope can only have these chars ["+r.Charset+"]",
		"invalid characters ["+invalidChars+"]",
	), false
}

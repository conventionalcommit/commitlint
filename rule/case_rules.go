package rule

import (
	"fmt"
	"slices"

	"github.com/conventionalcommit/commitlint/internal/casing"
	"github.com/conventionalcommit/commitlint/lint"
)

// Compile-time interface checks
var (
	_ lint.Rule = (*TypeCaseRule)(nil)
	_ lint.Rule = (*ScopeCaseRule)(nil)
	_ lint.Rule = (*DescriptionCaseRule)(nil)
	_ lint.Rule = (*BodyCaseRule)(nil)
	_ lint.Rule = (*HeaderCaseRule)(nil)
)

// CaseValues is the ordered list of supported case formats.
// Use the casing.* constants (e.g. casing.Lower, casing.Pascal) when
// referring to individual formats in code.
var CaseValues = casing.All

func caseIssue(field, caseFormat string) *lint.Issue {
	return lint.NewIssue(fmt.Sprintf("%s must be in %s", field, caseFormat))
}

func applyCaseArg(dst *string, ruleName string, setting lint.RuleSetting) error {
	if err := setStringArg(dst, setting.Argument); err != nil {
		return errInvalidArg(ruleName, err)
	}
	if slices.Contains(casing.All, *dst) {
		return nil
	}
	return errInvalidArg(ruleName, fmt.Errorf("unknown case %q, valid values: %v", *dst, casing.All))
}

// TypeCaseRule validates that commit type matches a given case format.
// Argument: one of the casing.* constants (e.g. casing.Lower).
type TypeCaseRule struct{ Case string }

func (r *TypeCaseRule) Name() string { return "type-case" }
func (r *TypeCaseRule) Apply(s lint.RuleSetting) error {
	return applyCaseArg(&r.Case, r.Name(), s)
}

func (r *TypeCaseRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if casing.Check(msg.Type(), r.Case) {
		return nil, true
	}
	return caseIssue("type", r.Case), false
}

// ScopeCaseRule validates that commit scope matches a given case format.
// An empty scope is always accepted (scope is optional by convention).
// Argument: one of the casing.* constants.
type ScopeCaseRule struct{ Case string }

func (r *ScopeCaseRule) Name() string { return "scope-case" }
func (r *ScopeCaseRule) Apply(s lint.RuleSetting) error {
	return applyCaseArg(&r.Case, r.Name(), s)
}

func (r *ScopeCaseRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	// scope is optional; skip the check when absent
	if msg.Scope() == "" {
		return nil, true
	}
	if casing.Check(msg.Scope(), r.Case) {
		return nil, true
	}
	return caseIssue("scope", r.Case), false
}

// DescriptionCaseRule validates that the commit description (subject) matches a
// given case format.
// Argument: one of the casing.* constants.
type DescriptionCaseRule struct{ Case string }

func (r *DescriptionCaseRule) Name() string { return "description-case" }
func (r *DescriptionCaseRule) Apply(s lint.RuleSetting) error {
	return applyCaseArg(&r.Case, r.Name(), s)
}

func (r *DescriptionCaseRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if casing.Check(msg.Description(), r.Case) {
		return nil, true
	}
	return caseIssue("description", r.Case), false
}

// BodyCaseRule validates that the commit body as a whole matches a given case
// format. An empty body always passes.
// Argument: one of the casing.* constants.
type BodyCaseRule struct{ Case string }

func (r *BodyCaseRule) Name() string { return "body-case" }
func (r *BodyCaseRule) Apply(s lint.RuleSetting) error {
	return applyCaseArg(&r.Case, r.Name(), s)
}

func (r *BodyCaseRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	body := msg.Body()
	if body == "" {
		return nil, true
	}
	if casing.Check(body, r.Case) {
		return nil, true
	}
	return caseIssue("body", r.Case), false
}

// HeaderCaseRule validates that the commit header matches a given case format.
// Argument: one of the casing.* constants.
type HeaderCaseRule struct{ Case string }

func (r *HeaderCaseRule) Name() string { return "header-case" }
func (r *HeaderCaseRule) Apply(s lint.RuleSetting) error {
	return applyCaseArg(&r.Case, r.Name(), s)
}

func (r *HeaderCaseRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if casing.Check(msg.Header(), r.Case) {
		return nil, true
	}
	return caseIssue("header", r.Case), false
}

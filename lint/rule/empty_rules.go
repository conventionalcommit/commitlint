package rule

import (
	"github.com/conventionalcommit/commitlint/lint"
)

// Compile-time interface checks
var (
	_ lint.Rule = (*TypeEmptyRule)(nil)
	_ lint.Rule = (*ScopeEmptyRule)(nil)
	_ lint.Rule = (*BodyEmptyRule)(nil)
	_ lint.Rule = (*FooterEmptyRule)(nil)
	_ lint.Rule = (*DescriptionEmptyRule)(nil)
)

// TypeEmptyRule validates that the commit type is not empty.
type TypeEmptyRule struct{}

func (r *TypeEmptyRule) Name() string                   { return "type-empty" }
func (r *TypeEmptyRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *TypeEmptyRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if msg.Type() != "" {
		return nil, true
	}
	return lint.NewIssue("type must not be empty"), false
}

// ScopeEmptyRule validates that the commit scope is not empty.
type ScopeEmptyRule struct{}

func (r *ScopeEmptyRule) Name() string                   { return "scope-empty" }
func (r *ScopeEmptyRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *ScopeEmptyRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if msg.Scope() != "" {
		return nil, true
	}
	return lint.NewIssue("scope must not be empty"), false
}

// BodyEmptyRule validates that the commit body is not empty.
type BodyEmptyRule struct{}

func (r *BodyEmptyRule) Name() string                   { return "body-empty" }
func (r *BodyEmptyRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *BodyEmptyRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if msg.Body() != "" {
		return nil, true
	}
	return lint.NewIssue("body must not be empty"), false
}

// FooterEmptyRule validates that the commit footer is not empty.
type FooterEmptyRule struct{}

func (r *FooterEmptyRule) Name() string                   { return "footer-empty" }
func (r *FooterEmptyRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *FooterEmptyRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if msg.Footer() != "" {
		return nil, true
	}
	return lint.NewIssue("footer must not be empty"), false
}

// DescriptionEmptyRule validates that the commit description is not empty.
type DescriptionEmptyRule struct{}

func (r *DescriptionEmptyRule) Name() string                   { return "description-empty" }
func (r *DescriptionEmptyRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *DescriptionEmptyRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if msg.Description() != "" {
		return nil, true
	}
	return lint.NewIssue("description must not be empty"), false
}

package rule

import (
	"github.com/conventionalcommit/commitlint/lint"
)

// Compile-time interface checks
var (
	_ lint.Rule = (*BodyLeadingBlankRule)(nil)
	_ lint.Rule = (*FooterLeadingBlankRule)(nil)
)

// BodyLeadingBlankRule checks that when a body is present, it begins with a blank line
// (i.e. the raw commit message has an empty line between header and body).
type BodyLeadingBlankRule struct{}

func (r *BodyLeadingBlankRule) Name() string                   { return "body-leading-blank" }
func (r *BodyLeadingBlankRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *BodyLeadingBlankRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if msg.Body() == "" || msg.BodyBlankLines() > 0 {
		return nil, true
	}
	return lint.NewIssue("body must have a leading blank line"), false
}

// FooterLeadingBlankRule checks that when a footer is present, it begins with a blank line
// (i.e. there is an empty line between body/header and footer).
type FooterLeadingBlankRule struct{}

func (r *FooterLeadingBlankRule) Name() string                   { return "footer-leading-blank" }
func (r *FooterLeadingBlankRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *FooterLeadingBlankRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if msg.Footer() == "" || msg.FooterBlankLines() > 0 {
		return nil, true
	}
	return lint.NewIssue("footer must have a leading blank line"), false
}

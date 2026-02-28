package rule

import (
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// Compile-time interface checks
var (
	_ lint.Rule = (*BodyLeadingBlankRule)(nil)
	_ lint.Rule = (*FooterLeadingBlankRule)(nil)
)

// BodyLeadingBlankRule checks that when a body is present, it begins with a blank line
// (i.e. the raw commit message has an empty line between header and body).
// Per the Conventional Commits spec the parser already trims the leading blank, but we
// can detect its presence via the full commit Message(): the body section should start
// after two newlines (header + blank line + body).
type BodyLeadingBlankRule struct{}

func (r *BodyLeadingBlankRule) Name() string                   { return "body-leading-blank" }
func (r *BodyLeadingBlankRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *BodyLeadingBlankRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	body := msg.Body()
	if body == "" {
		return nil, true
	}
	// The full message should have "\n\n" separating header from body.
	raw := msg.Message()
	headerEnd := strings.Index(raw, "\n")
	if headerEnd == -1 {
		// No newline at all, body is non-empty but raw message has no newline;
		// treat as no leading blank.
		return lint.NewIssue("body must have a leading blank line"), false
	}
	rest := raw[headerEnd:]
	if strings.HasPrefix(rest, "\n\n") {
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
	footer := msg.Footer()
	if footer == "" {
		return nil, true
	}
	raw := msg.Message()
	// Footer should be preceded by "\n\n"
	footerIdx := strings.LastIndex(raw, footer)
	if footerIdx < 2 {
		return lint.NewIssue("footer must have a leading blank line"), false
	}
	if raw[footerIdx-2:footerIdx] == "\n\n" {
		return nil, true
	}
	return lint.NewIssue("footer must have a leading blank line"), false
}

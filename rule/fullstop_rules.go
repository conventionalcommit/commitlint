package rule

import (
	"fmt"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// Compile-time interface checks
var (
	_ lint.Rule = (*HeaderFullStopRule)(nil)
	_ lint.Rule = (*BodyFullStopRule)(nil)
	_ lint.Rule = (*DescriptionFullStopRule)(nil)
)

// applyFullStopArg extracts a single character string from setting.
func applyFullStopArg(dst *string, ruleName string, setting lint.RuleSetting) error {
	if err := setStringArg(dst, setting.Argument); err != nil {
		return errInvalidArg(ruleName, err)
	}
	return nil
}

// HeaderFullStopRule checks that the header does NOT end with a given character.
// Default character is ".".
type HeaderFullStopRule struct{ Char string }

func (r *HeaderFullStopRule) Name() string { return "header-full-stop" }
func (r *HeaderFullStopRule) Apply(s lint.RuleSetting) error {
	return applyFullStopArg(&r.Char, r.Name(), s)
}

func (r *HeaderFullStopRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	if !strings.HasSuffix(msg.Header(), r.Char) {
		return nil, true
	}
	return lint.NewIssue(fmt.Sprintf("header must not end with %q", r.Char)), false
}

// BodyFullStopRule checks that the body does NOT end with a given character.
type BodyFullStopRule struct{ Char string }

func (r *BodyFullStopRule) Name() string { return "body-full-stop" }
func (r *BodyFullStopRule) Apply(s lint.RuleSetting) error {
	return applyFullStopArg(&r.Char, r.Name(), s)
}

func (r *BodyFullStopRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	body := msg.Body()
	if body == "" {
		return nil, true
	}
	if !strings.HasSuffix(body, r.Char) {
		return nil, true
	}
	return lint.NewIssue(fmt.Sprintf("body must not end with %q", r.Char)), false
}

// DescriptionFullStopRule checks that the description does NOT end with a given character.
type DescriptionFullStopRule struct{ Char string }

func (r *DescriptionFullStopRule) Name() string { return "description-full-stop" }
func (r *DescriptionFullStopRule) Apply(s lint.RuleSetting) error {
	return applyFullStopArg(&r.Char, r.Name(), s)
}

func (r *DescriptionFullStopRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	desc := msg.Description()
	if desc == "" {
		return nil, true
	}
	if !strings.HasSuffix(desc, r.Char) {
		return nil, true
	}
	return lint.NewIssue(fmt.Sprintf("description must not end with %q", r.Char)), false
}

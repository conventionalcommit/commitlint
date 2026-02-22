package rule

import (
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// Compile-time interface checks.
var (
	_ lint.Rule = (*SignedOffByRule)(nil)
	_ lint.Rule = (*TrailerExistsRule)(nil)
)

// SignedOffByRule checks that at least one footer note token matches the
// configured value (default "Signed-off-by").
type SignedOffByRule struct{ Value string }

func (r *SignedOffByRule) Name() string { return "signed-off-by" }
func (r *SignedOffByRule) Apply(s lint.RuleSetting) error {
	if err := applyStringArg(&r.Value, r.Name(), s); err != nil {
		return err
	}
	// Normalize: strip a trailing ":" so both "Signed-off-by:" and
	// "Signed-off-by" match note tokens returned by the parser.
	r.Value = strings.TrimSuffix(strings.TrimSpace(r.Value), ":")
	return nil
}

func (r *SignedOffByRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	for _, note := range msg.Notes() {
		if note.Token() == r.Value {
			return nil, true
		}
	}
	return lint.NewIssue("message must contain trailer " + r.Value), false
}

// TrailerExistsRule checks that at least one footer note has a token matching
// the configured value.
//
// This is a generalized version of signed-off-by: any trailer token can be
// required.
type TrailerExistsRule struct{ Value string }

func (r *TrailerExistsRule) Name() string { return "trailer-exists" }
func (r *TrailerExistsRule) Apply(s lint.RuleSetting) error {
	if err := applyStringArg(&r.Value, r.Name(), s); err != nil {
		return err
	}
	// Normalize: strip a trailing ":" so both "Co-authored-by:" and
	// "Co-authored-by" match note tokens returned by the parser.
	r.Value = strings.TrimSuffix(strings.TrimSpace(r.Value), ":")
	return nil
}

func (r *TrailerExistsRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	for _, note := range msg.Notes() {
		if note.Token() == r.Value {
			return nil, true
		}
	}
	return lint.NewIssue("message must contain trailer " + r.Value), false
}

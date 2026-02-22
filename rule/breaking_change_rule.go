package rule

import "github.com/conventionalcommit/commitlint/lint"

// Compile-time interface check.
var _ lint.Rule = (*BreakingChangeExclamationMarkRule)(nil)

// BreakingChangeExclamationMarkRule implements the breaking-change-exclamation-mark rule.
//
// It enforces that breaking changes are signalled consistently via an XNOR:
//   - PASS when the header contains "!" AND a footer note has token
//     "BREAKING CHANGE" or "BREAKING-CHANGE".
//   - PASS when neither signal is present.
//   - FAIL when exactly one of the two is present.
//
// Detection uses the parsed Notes() (footer tokens), not raw string scanning.
type BreakingChangeExclamationMarkRule struct{}

func (r *BreakingChangeExclamationMarkRule) Name() string                   { return "breaking-change-exclamation-mark" }
func (r *BreakingChangeExclamationMarkRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *BreakingChangeExclamationMarkRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	hasExclamation := msg.IsBreakingChange()

	// Detect "BREAKING CHANGE" or "BREAKING-CHANGE" via parsed footer notes.
	hasFooterBreaking := false
	for _, note := range msg.Notes() {
		t := note.Token()
		if t == "BREAKING CHANGE" || t == "BREAKING-CHANGE" {
			hasFooterBreaking = true
			break
		}
	}

	// XNOR: pass when both present or both absent.
	if hasExclamation == hasFooterBreaking {
		return nil, true
	}

	if hasExclamation && !hasFooterBreaking {
		return lint.NewIssue(
			"breaking change exclamation mark in header requires BREAKING CHANGE in footer",
		), false
	}
	return lint.NewIssue(
		"BREAKING CHANGE in footer requires exclamation mark in header",
	), false
}

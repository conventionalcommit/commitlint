package rule

import (
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// Compile-time interface checks
var _ lint.Rule = (*HeaderTrimRule)(nil)

// HeaderTrimRule checks that the header has no leading or trailing whitespace.
type HeaderTrimRule struct{}

func (r *HeaderTrimRule) Name() string                   { return "header-trim" }
func (r *HeaderTrimRule) Apply(_ lint.RuleSetting) error { return nil }
func (r *HeaderTrimRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	h := msg.Header()
	if h == strings.TrimSpace(h) {
		return nil, true
	}
	return lint.NewIssue("header must not have leading or trailing whitespace"), false
}

package rule

import (
	"fmt"
	"sort"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

var _ lint.Rule = (*FooterEnumRule)(nil)

// FooterEnumRule to validate footer tokens
type FooterEnumRule struct {
	Tokens []string
}

// Name return name of the rule
func (r *FooterEnumRule) Name() string { return "footer-enum" }

// Apply sets the needed argument for the rule
func (r *FooterEnumRule) Apply(setting lint.RuleSetting) error {
	err := setStringArrArg(&r.Tokens, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	// sorting the string elements for binary search
	sort.Strings(r.Tokens)
	return nil
}

// Validate validates FooterEnumRule
func (r *FooterEnumRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	var invalids []string

	for _, note := range msg.Notes() {
		isFound := search(r.Tokens, note.Token())
		if !isFound {
			invalids = append(invalids, note.Token())
		}
	}

	if len(invalids) == 0 {
		return nil, true
	}

	desc := fmt.Sprintf("you can use one of %v", r.Tokens)
	info := fmt.Sprintf("[%s] tokens are not allowed", strings.Join(invalids, ", "))
	return lint.NewIssue(desc, info), false
}

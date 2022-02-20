package rule

import (
	"fmt"
	"sort"

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
func (r *FooterEnumRule) Validate(msg lint.Commit) ([]string, bool) {
	msgs := []string{}

	for _, note := range msg.Notes() {
		isFound := search(r.Tokens, note.Token())
		if !isFound {
			errMsg := fmt.Sprintf("footer token '%s' is not allowed, you can use one of %v", note.Token(), r.Tokens)
			msgs = append(msgs, errMsg)
		}
	}

	if len(msgs) == 0 {
		return nil, true
	}

	return msgs, false
}

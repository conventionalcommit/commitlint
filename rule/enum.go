package rule

import (
	"fmt"
	"sort"

	"github.com/conventionalcommit/commitlint/lint"
)

// ScopeEnumRule to validate max length of header
type ScopeEnumRule struct {
	Scopes []string

	AllowEmpty bool
}

// Name return name of the rule
func (r *ScopeEnumRule) Name() string { return "scope-enum" }

// Validate validates ScopeEnumRule
func (r *ScopeEnumRule) Validate(msg *lint.Commit) ([]string, bool) {
	if msg.Scope() == "" {
		if r.AllowEmpty {
			return nil, true
		}
		errMsg := fmt.Sprintf("empty scope is not allowed, you can use one of %v", r.Scopes)
		return []string{errMsg}, false
	}

	isFound := search(r.Scopes, msg.Scope())
	if !isFound {
		errMsg := fmt.Sprintf("scope '%s' is not allowed, you can use one of %v", msg.Scope(), r.Scopes)
		return []string{errMsg}, false
	}
	return nil, true
}

// Apply sets the needed argument for the rule
func (r *ScopeEnumRule) Apply(setting lint.RuleSetting) error {
	err := setStringArrArg(&r.Scopes, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}

	allowEmpty, ok := setting.Flags["allow-empty"]
	if ok {
		err := setBoolArg(&r.AllowEmpty, allowEmpty)
		if err != nil {
			return errInvalidFlag(r.Name(), "allow-empty", err)
		}
	}

	// sorting the string elements for binary search
	sort.Strings(r.Scopes)
	return nil
}

// TypeEnumRule to validate types
type TypeEnumRule struct {
	Types []string
}

// Name return name of the rule
func (r *TypeEnumRule) Name() string { return "type-enum" }

// Validate validates TypeEnumRule
func (r *TypeEnumRule) Validate(msg *lint.Commit) ([]string, bool) {
	isFound := search(r.Types, msg.Type())
	if !isFound {
		errMsg := fmt.Sprintf("type '%s' is not allowed, you can use one of %v", msg.Type(), r.Types)
		return []string{errMsg}, false
	}
	return nil, true
}

// Apply sets the needed argument for the rule
func (r *TypeEnumRule) Apply(setting lint.RuleSetting) error {
	err := setStringArrArg(&r.Types, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	// sorting the string elements for binary search
	sort.Strings(r.Types)
	return nil
}

// FooterEnumRule to validate footer tokens
type FooterEnumRule struct {
	Tokens []string
}

// Name return name of the rule
func (r *FooterEnumRule) Name() string { return "footer-enum" }

// Validate validates FooterEnumRule
func (r *FooterEnumRule) Validate(msg *lint.Commit) ([]string, bool) {
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

func search(arr []string, toFind string) bool {
	ind := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= toFind
	})
	return ind < len(arr) && arr[ind] == toFind
}

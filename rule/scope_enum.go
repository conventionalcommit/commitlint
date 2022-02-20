package rule

import (
	"fmt"
	"sort"

	"github.com/conventionalcommit/commitlint/lint"
)

var _ lint.Rule = (*ScopeEnumRule)(nil)

// ScopeEnumRule to validate max length of header
type ScopeEnumRule struct {
	Scopes []string

	AllowEmpty bool
}

// Name return name of the rule
func (r *ScopeEnumRule) Name() string { return "scope-enum" }

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

// Validate validates ScopeEnumRule
func (r *ScopeEnumRule) Validate(msg lint.Commit) ([]string, bool) {
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

package rule

import (
	"fmt"
	"sort"

	"github.com/conventionalcommit/commitlint/message"
)

// ScopeEnumRule to validate max length of header
type ScopeEnumRule struct {
	Scopes []string
}

// Name return name of the rule
func (r *ScopeEnumRule) Name() string { return "scope-enum" }

// Validate validates ScopeEnumRule
func (r *ScopeEnumRule) Validate(msg *message.Commit) (string, bool) {
	isFound := search(r.Scopes, msg.Header.Scope)
	if !isFound {
		errMsg := fmt.Sprintf("scope '%s' is not allowed, you can use one of %v", msg.Header.Scope, r.Scopes)
		return errMsg, false
	}
	return "", true
}

// Apply sets the needed argument for the rule
func (r *ScopeEnumRule) Apply(arg interface{}, flags map[string]interface{}) error {
	err := setStringArrArg(&r.Scopes, arg, r.Name())
	if err != nil {
		return err
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
func (r *TypeEnumRule) Validate(msg *message.Commit) (string, bool) {
	isFound := search(r.Types, msg.Header.Type)
	if !isFound {
		errMsg := fmt.Sprintf("type '%s' is not allowed, you can use one of %v", msg.Header.Type, r.Types)
		return errMsg, false
	}
	return "", true
}

// Apply sets the needed argument for the rule
func (r *TypeEnumRule) Apply(arg interface{}, flags map[string]interface{}) error {
	err := setStringArrArg(&r.Types, arg, r.Name())
	if err != nil {
		return err
	}
	// sorting the string elements for binary search
	sort.Strings(r.Types)
	return nil
}

func search(arr []string, toFind string) bool {
	ind := sort.Search(len(arr), func(i int) bool {
		return toFind <= arr[i]
	})
	return ind < len(arr) && arr[ind] == toFind
}

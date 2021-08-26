package rule

import (
	"fmt"

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

// SetAndCheckArgument sets the needed argument for the rule
func (r *ScopeEnumRule) SetAndCheckArgument(arg interface{}) error {
	return setStringArrArg(&r.Scopes, arg, r.Name())
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

// SetAndCheckArgument sets the needed argument for the rule
func (r *TypeEnumRule) SetAndCheckArgument(arg interface{}) error {
	return setStringArrArg(&r.Types, arg, r.Name())
}

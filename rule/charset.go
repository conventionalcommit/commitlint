package rule

import (
	"fmt"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// ScopeCharsetRule to validate max length of header
type ScopeCharsetRule struct {
	Charset string
}

// Name return name of the rule
func (r *ScopeCharsetRule) Name() string { return "scope-charset" }

// Validate validates ScopeCharsetRule
func (r *ScopeCharsetRule) Validate(msg *lint.Commit) (string, bool) {
	invalidChar, isValid := checkCharset(r.Charset, msg.Header.Scope)
	if !isValid {
		errMsg := fmt.Sprintf("scope contains invalid char '%s', allowed chars are [%s]", invalidChar, r.Charset)
		return errMsg, false
	}
	return "", true
}

// Apply sets the needed argument for the rule
func (r *ScopeCharsetRule) Apply(arg interface{}, flags map[string]interface{}) error {
	err := setStringArg(&r.Charset, arg)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// TypeCharsetRule to validate max length of header
type TypeCharsetRule struct {
	Charset string
}

// Name return name of the rule
func (r *TypeCharsetRule) Name() string { return "type-charset" }

// Validate validates TypeCharsetRule
func (r *TypeCharsetRule) Validate(msg *lint.Commit) (string, bool) {
	invalidChar, isValid := checkCharset(r.Charset, msg.Header.Type)
	if !isValid {
		errMsg := fmt.Sprintf("type contains invalid char '%s', allowed chars are [%s]", invalidChar, r.Charset)
		return errMsg, false
	}
	return "", true
}

// Apply sets the needed argument for the rule
func (r *TypeCharsetRule) Apply(arg interface{}, flags map[string]interface{}) error {
	err := setStringArg(&r.Charset, arg)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

func checkCharset(charset, toCheck string) (string, bool) {
	for _, ch := range toCheck {
		if !strings.ContainsRune(charset, ch) {
			return string(ch), false
		}
	}
	return "", true
}

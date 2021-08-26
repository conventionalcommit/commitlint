package rule

import (
	"fmt"
	"strings"

	"github.com/conventionalcommit/commitlint/message"
)

// ScopeCharsetRule to validate max length of header
type ScopeCharsetRule struct {
	Charset string
}

// Name return name of the rule
func (r *ScopeCharsetRule) Name() string { return "scope-charset" }

// Validate validates ScopeCharsetRule
func (r *ScopeCharsetRule) Validate(msg *message.Commit) (string, bool) {
	invalidChar, isValid := checkCharset(r.Charset, msg.Header.Scope)
	if !isValid {
		errMsg := fmt.Sprintf("scope contains invalid char '%s', allowed chars are [%s]", invalidChar, r.Charset)
		return errMsg, false
	}
	return "", true
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *ScopeCharsetRule) SetAndCheckArgument(arg interface{}) error {
	return setStringArg(&r.Charset, arg, r.Name())
}

// TypeCharsetRule to validate max length of header
type TypeCharsetRule struct {
	Charset string
}

// Name return name of the rule
func (r *TypeCharsetRule) Name() string { return "type-charset" }

// Validate validates TypeCharsetRule
func (r *TypeCharsetRule) Validate(msg *message.Commit) (string, bool) {
	invalidChar, isValid := checkCharset(r.Charset, msg.Header.Type)
	if !isValid {
		errMsg := fmt.Sprintf("type contains invalid char '%s', allowed chars are [%s]", invalidChar, r.Charset)
		return errMsg, false
	}
	return "", true
}

// SetAndCheckArgument sets the needed argument for the rule
func (r *TypeCharsetRule) SetAndCheckArgument(arg interface{}) error {
	return setStringArg(&r.Charset, arg, r.Name())
}

func checkCharset(charset, toCheck string) (string, bool) {
	for _, ch := range toCheck {
		if !strings.ContainsRune(charset, ch) {
			return string(ch), false
		}
	}
	return "", true
}

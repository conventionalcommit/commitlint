package rule

import (
	"fmt"

	"github.com/conventionalcommit/commitlint/lint"
)

// HeadMinLenRule to validate min length of header
type HeadMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *HeadMinLenRule) Name() string { return "header-min-length" }

// Validate validates HeadMinLenRule
func (r *HeadMinLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMinLen(r.CheckLen, msg.Header.FullHeader)
}

// Apply sets the needed argument for the rule
func (r *HeadMinLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// BodyMinLenRule to validate min length of body
type BodyMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *BodyMinLenRule) Name() string { return "body-min-length" }

// Validate validates BodyMinLenRule
func (r *BodyMinLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMinLen(r.CheckLen, msg.Body)
}

// Apply sets the needed argument for the rule
func (r *BodyMinLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// FooterMinLenRule to validate min length of footer
type FooterMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *FooterMinLenRule) Name() string { return "footer-min-length" }

// Validate validates FooterMinLenRule
func (r *FooterMinLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMinLen(r.CheckLen, msg.Footer.FullFooter)
}

// Apply sets the needed argument for the rule
func (r *FooterMinLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// TypeMinLenRule to validate min length of type
type TypeMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *TypeMinLenRule) Name() string { return "type-min-length" }

// Validate validates TypeMinLenRule
func (r *TypeMinLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMinLen(r.CheckLen, msg.Header.Type)
}

// Apply sets the needed argument for the rule
func (r *TypeMinLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// ScopeMinLenRule to validate min length of scope
type ScopeMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *ScopeMinLenRule) Name() string { return "scope-min-length" }

// Validate validates ScopeMinLenRule
func (r *ScopeMinLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMinLen(r.CheckLen, msg.Header.Scope)
}

// Apply sets the needed argument for the rule
func (r *ScopeMinLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

// DescriptionMinLenRule to validate min length of description
type DescriptionMinLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *DescriptionMinLenRule) Name() string { return "description-min-length" }

// Validate validates DescriptionMinLenRule
func (r *DescriptionMinLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMinLen(r.CheckLen, msg.Header.Description)
}

// Apply sets the needed argument for the rule
func (r *DescriptionMinLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	return setIntArg(&r.CheckLen, arg, r.Name())
}

func checkMinLen(checkLen int, toCheck string) (string, bool) {
	actualLen := len(toCheck)
	if actualLen < checkLen {
		errMsg := fmt.Sprintf("length is %d, should have atleast %d chars", actualLen, checkLen)
		return errMsg, false
	}
	return "", true
}

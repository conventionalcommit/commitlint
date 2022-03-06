package rule

import (
	"github.com/conventionalcommit/commitlint/lint"
)

var _ lint.Rule = (*BodyMaxLineLenRule)(nil)

// BodyMaxLineLenRule to validate max line length of body
type BodyMaxLineLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *BodyMaxLineLenRule) Name() string { return "body-max-line-length" }

// Apply sets the needed argument for the rule
func (r *BodyMaxLineLenRule) Apply(setting lint.RuleSetting) error {
	err := setIntArg(&r.CheckLen, setting.Argument)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// Validate validates BodyMaxLineLenRule rule
func (r *BodyMaxLineLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLineLength("body", r.CheckLen, msg.Body())
}

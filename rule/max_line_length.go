package rule

import (
	"fmt"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// BodyMaxLineLenRule to validate max line length of body
type BodyMaxLineLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *BodyMaxLineLenRule) Name() string { return "body-max-line-length" }

// Validate validates BodyMaxLineLenRule rule
func (r *BodyMaxLineLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMaxLineLength(r.CheckLen, msg.Body)
}

// Apply sets the needed argument for the rule
func (r *BodyMaxLineLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	err := setIntArg(&r.CheckLen, arg)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

// FooterMaxLineLenRule to validate max line length of footer
type FooterMaxLineLenRule struct {
	CheckLen int
}

// Name return name of the rule
func (r *FooterMaxLineLenRule) Name() string { return "footer-max-line-length" }

// Validate validates FooterMaxLineLenRule rule
func (r *FooterMaxLineLenRule) Validate(msg *lint.Commit) (string, bool) {
	return checkMaxLineLength(r.CheckLen, msg.Footer.FullFooter)
}

// Apply sets the needed argument for the rule
func (r *FooterMaxLineLenRule) Apply(arg interface{}, flags map[string]interface{}) error {
	err := setIntArg(&r.CheckLen, arg)
	if err != nil {
		return errInvalidArg(r.Name(), err)
	}
	return nil
}

func checkMaxLineLength(checkLen int, toCheck string) (string, bool) {
	lines := strings.Split(toCheck, "\n")
	for index, line := range lines {
		actualLen := len(line)
		if actualLen > checkLen {
			errMsg := fmt.Sprintf("in line %d, length is %d, should have less than %d chars", index+1, actualLen, checkLen)
			return errMsg, false
		}
	}
	return "", true
}

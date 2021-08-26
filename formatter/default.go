// Package formatter contains lint result formatters
package formatter

import (
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// DefaultFormatter represent default formatter
type DefaultFormatter struct{}

// Name returns name of formatter
func (f *DefaultFormatter) Name() string { return "default" }

// Format formats the lint.Result
func (f *DefaultFormatter) Format(res *lint.Result) (string, error) {
	msg := formatResult(res)
	return strings.Trim(msg, "\n"), nil
}

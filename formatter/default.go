// Package formatter contains lint result formatters
package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// DefaultFormatter represent default formatter
type DefaultFormatter struct{}

// Name returns name of formatter
func (f *DefaultFormatter) Name() string { return "default" }

// Format formats the lint.Failure
func (f *DefaultFormatter) Format(res *lint.Failure) (string, error) {
	return formatFailure(res), nil
}

func formatFailure(res *lint.Failure) string {
	if res.IsOK() {
		return " ✔ commit message"
	}
	return writeFailure(res)
}

func writeFailure(res *lint.Failure) string {
	str := &strings.Builder{}

	str.WriteString("commitlint")
	fmt.Fprintf(str, "\n\n→ input: %s", strconv.Quote(truncate(25, res.Input())))

	errs, warns, others := bySeverity(res)

	writeRuleFailure(str, "Errors", errs, "❌")
	writeRuleFailure(str, "Warnings", warns, "!")
	writeRuleFailure(str, "Other Severities", others, "?")

	fmt.Fprintf(str, "\n\nTotal %d errors, %d warnings, %d other severities", len(errs), len(warns), len(others))
	return str.String()
}

func writeRuleFailure(w *strings.Builder, title string, resArr []*lint.RuleFailure, sign string) {
	if len(resArr) == 0 {
		return
	}

	fmt.Fprint(w, "\n\n"+title+":")
	for _, msg := range resArr {
		fmt.Fprintf(w, "\n    %s %s: %s", sign, msg.RuleName(), msg.Message())
	}
}

func truncate(maxSize int, input string) string {
	if len(input) < maxSize {
		return input
	}
	return input[:maxSize-3] + "..."
}

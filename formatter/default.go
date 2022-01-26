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

	quotedStr := strconv.Quote(truncate(25, res.Input()))

	str.WriteString("commitlint\n")
	str.WriteString("\n→ input: " + quotedStr)

	errs, warns, others := bySeverity(res)

	writeRuleFailure(str, "❌", "Errors", errs)
	writeRuleFailure(str, "!", "Warnings", warns)
	writeRuleFailure(str, "?", "Other Severities", others)

	fmt.Fprintf(str, "\n\nTotal %d errors, %d warnings, %d other severities", len(errs), len(warns), len(others))
	return strings.Trim(str.String(), "\n")
}

func writeRuleFailure(w *strings.Builder, sign, title string, resArr []*lint.RuleFailure) {
	if len(resArr) == 0 {
		return
	}

	w.WriteString("\n\n" + title + ":")
	for _, ruleRes := range resArr {
		writeMessages(w, ruleRes, sign)
	}
}

func writeMessages(w *strings.Builder, ruleRes *lint.RuleFailure, sign string) {
	msgs := ruleRes.Message()

	if len(msgs) == 0 {
		return
	}

	space := "  "

	if len(msgs) == 1 {
		msg := msgs[0]
		// ❌ rule-name: message
		fmt.Fprintf(w, "\n%s %s: %s", space+sign, ruleRes.Name(), msg)
		return
	}

	// ❌ rule-name:
	//    - message1
	//    - message2
	fmt.Fprintf(w, "\n%s %s:", space+sign, ruleRes.Name())
	for _, msg := range ruleRes.Message() {
		fmt.Fprintf(w, "\n%s - %s", space+space, msg)
	}
}

func truncate(maxSize int, input string) string {
	if len(input) < maxSize {
		return input
	}
	return input[:maxSize-3] + "..."
}

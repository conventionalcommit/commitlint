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

// Format formats the lint.Result
func (f *DefaultFormatter) Format(res *lint.Result) (string, error) {
	return formatResult(res), nil
}

func formatResult(res *lint.Result) string {
	if res.IsOK() {
		return " ✔ commit message"
	}
	return writeResult(res)
}

func writeResult(res *lint.Result) string {
	str := &strings.Builder{}

	str.WriteString("commitlint")
	fmt.Fprintf(str, "\n\n→ input: %s", strconv.Quote(truncate(25, res.Input())))

	if res.HasErrors() {
		writeLintResult(str, "Errors", res.Errors(), "❌")
	}

	if res.HasWarns() {
		writeLintResult(str, "Warnings", res.Warns(), "!")
	}

	fmt.Fprintf(str, "\n\nTotal %d errors, %d warnings", len(res.Errors()), len(res.Warns()))
	return str.String()
}

func writeLintResult(w *strings.Builder, title string, resArr []lint.RuleResult, sign string) {
	fmt.Fprint(w, "\n\n"+title+":")
	for _, msg := range resArr {
		fmt.Fprintf(w, "\n    %s %s: %s", sign, msg.Name, msg.Message)
	}
}

func truncate(maxSize int, input string) string {
	if len(input) < maxSize {
		return input
	}
	return input[:maxSize-3] + "..."
}

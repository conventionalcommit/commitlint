// Package formatter contains lint result formatters
package formatter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

const (
	truncateSize = 25
)

// DefaultFormatter represent default formatter
type DefaultFormatter struct{}

// Name returns name of formatter
func (f *DefaultFormatter) Name() string { return "default" }

// Format formats the lint.Failure
func (f *DefaultFormatter) Format(result *lint.Result) (string, error) {
	return f.formatFailure(result.Input(), result.Issues()), nil
}

func (f *DefaultFormatter) formatFailure(msg string, issues []*lint.Issue) string {
	if len(issues) == 0 {
		return " ✔ commit message"
	}
	return f.writeFailure(msg, issues)
}

func (f *DefaultFormatter) writeFailure(msg string, issues []*lint.Issue) string {
	str := &strings.Builder{}

	quotedStr := strconv.Quote(truncate(truncateSize, msg))

	str.WriteString("commitlint\n")
	str.WriteString("\n→ input: " + quotedStr)

	errs, warns, others := f.bySeverity(issues)

	f.writeIssues(str, "❌", "Errors", errs)
	f.writeIssues(str, "!", "Warnings", warns)
	f.writeIssues(str, "?", "Other Severities", others)

	fmt.Fprintf(str, "\n\nTotal %d errors, %d warnings, %d other severities", len(errs), len(warns), len(others))
	return strings.Trim(str.String(), "\n")
}

func (f *DefaultFormatter) writeIssues(w *strings.Builder, sign, title string, issues []*lint.Issue) {
	if len(issues) == 0 {
		return
	}

	w.WriteString("\n\n" + title + ":")
	for _, issue := range issues {
		f.writeIssue(w, sign, issue)
	}
}

func (f *DefaultFormatter) writeIssue(w *strings.Builder, sign string, issue *lint.Issue) {
	space := "  "

	// ❌ rule-name: description
	//    - info1
	//    - info2

	fmt.Fprintf(w, "\n%s %s: %s", space+sign, issue.RuleName(), issue.Description())
	for _, msg := range issue.Infos() {
		fmt.Fprintf(w, "\n%s - %s", space+space, msg)
	}
}

// bySeverity returns all messages with given severity
func (f *DefaultFormatter) bySeverity(issues []*lint.Issue) (errs, warns, others []*lint.Issue) {
	for _, r := range issues {
		switch r.Severity() {
		case lint.SeverityError:
			errs = append(errs, r)
		case lint.SeverityWarn:
			warns = append(warns, r)
		default:
			others = append(others, r)
		}
	}
	return errs, warns, others
}

func truncate(maxSize int, input string) string {
	if len(input) < maxSize {
		return input
	}
	return input[:maxSize-3] + "..."
}

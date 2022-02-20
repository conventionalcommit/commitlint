package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// JSONFormatter represent default formatter
type JSONFormatter struct{}

// Name returns name of formatter
func (f *JSONFormatter) Name() string { return "json" }

// Format formats the lint.Result
func (f *JSONFormatter) Format(result *lint.Result) (string, error) {
	output := make(map[string]interface{}, 4)

	output["input"] = result.Input()
	output["valid"] = result.IsOK()
	output["issues"] = f.formatIssue(result.Issues())

	formatted, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("json formatting failed: %w", err)
	}
	return strings.Trim(string(formatted), "\n"), nil
}

func (f *JSONFormatter) formatIssue(issues []*lint.Issue) []interface{} {
	formattedIssues := make([]interface{}, 0, len(issues))

	for _, issue := range issues {
		output := make(map[string]interface{})

		output["name"] = issue.Name()
		output["severity"] = issue.Severity()

		if len(issue.Message()) > 0 {
			output["messages"] = issue.Message()
		}

		formattedIssues = append(formattedIssues, output)
	}

	return formattedIssues
}

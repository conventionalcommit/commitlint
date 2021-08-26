package formatter

import (
	"encoding/json"
	"strings"

	"github.com/conventionalcommit/commitlint/lint"
)

// JSONFormatter represent default formatter
type JSONFormatter struct{}

// Name returns name of formatter
func (f *JSONFormatter) Name() string { return "json" }

// Format formats the lint.Result
func (f *JSONFormatter) Format(res *lint.Result) (string, error) {
	output := make(map[string]interface{}, 4)

	output["input"] = res.Input()
	output["status"] = res.IsOK()
	output["errors"] = f.formRuleResult(res.Errors())
	output["warnings"] = f.formRuleResult(res.Warns())

	msg, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(msg), "\n"), nil
}

func (f *JSONFormatter) formRuleResult(res []lint.RuleResult) []map[string]interface{} {
	outs := make([]map[string]interface{}, 0, len(res))

	for _, r := range res {
		output := make(map[string]interface{})

		output["name"] = r.Name
		output["message"] = r.Message

		outs = append(outs, output)
	}

	return outs
}

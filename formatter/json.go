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
func (f *JSONFormatter) Format(res *lint.Failure) (string, error) {
	output := make(map[string]interface{}, 4)

	errs, warns, others := bySeverity(res)

	output["input"] = res.Input()
	output["status"] = res.IsOK()

	output["errors"] = f.formRuleFailure(errs, false)
	output["warnings"] = f.formRuleFailure(warns, false)
	output["others"] = f.formRuleFailure(others, true)

	msg, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(msg), "\n"), nil
}

func (f *JSONFormatter) formRuleFailure(res []*lint.RuleFailure, includeSev bool) []map[string]interface{} {
	outs := make([]map[string]interface{}, 0, len(res))

	for _, r := range res {
		output := make(map[string]interface{})

		output["name"] = r.RuleName()
		output["message"] = r.Message()

		if includeSev {
			output["severity"] = r.Severity()
		}

		outs = append(outs, output)
	}

	return outs
}

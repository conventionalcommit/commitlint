package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

func TestLintRecoveredIssues(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"type-enum", "body-leading-blank", "body-max-line-length", "header-trim"}
	conf.Settings["body-max-line-length"] = lint.RuleSetting{Argument: 10}
	conf.Severity.Default = lint.SeverityWarn
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatal(err)
	}
	input := " unknown:expose records\n" + strings.Repeat("x", 11)
	result, err := linter.ParseAndLint(input)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"parser", "parser", "parser", "type-enum", "body-leading-blank", "body-max-line-length", "header-trim"}
	var rules []string
	for _, issue := range result.Issues() {
		rules = append(rules, issue.RuleName())
		severity := lint.SeverityWarn
		if issue.RuleName() == "parser" {
			severity = lint.SeverityError
		}
		if issue.Severity() != severity {
			t.Fatalf("severity=%s for %s", issue.Severity(), issue.RuleName())
		}
	}
	if !reflect.DeepEqual(rules, want) || result.Input() != input {
		t.Fatalf("rules=%v, want %v", rules, want)
	}
}

func TestLintAmbiguousHeaderSkipsRules(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"scope-empty", "body-empty", "footer-empty"}
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatal(err)
	}
	for _, input := range []string{"invalid header", "feat(store: expose records", "feat: expose records\xff"} {
		result, err := linter.ParseAndLint(input)
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Issues()) != 1 || result.Issues()[0].RuleName() != "parser" || result.Input() != input {
			t.Fatalf("ambiguous input %q produced guessed rule failures", input)
		}
	}
}

func TestLintRetainsFootersAfterInvalidValue(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"trailer-exists"}
	conf.Settings["trailer-exists"] = lint.RuleSetting{Argument: "Reviewed-by"}
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatal(err)
	}
	input := "feat: expose records\n\nRefs #9\nEmpty: \nReviewed-by: Sam"
	result, err := linter.ParseAndLint(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Issues()) != 1 || result.Issues()[0].RuleName() != "parser" {
		t.Fatalf("later valid footer was lost: %v", result.Issues())
	}
}

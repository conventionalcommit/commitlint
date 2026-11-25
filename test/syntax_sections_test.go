package test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/lint/rule"
)

func TestSyntaxLeadingBlankLines(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"body-leading-blank", "footer-leading-blank"}
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ name, line, blank string }{
		{"lf", "\n", "\n"},
		{"crlf", "\r\n", "\r\n"},
		{"lf then crlf", "\n", "\r\n"},
		{"crlf then lf", "\r\n", "\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := "feat: expose records" + tt.line + tt.blank + "Body text." + tt.line + tt.blank + "Refs #9" + tt.line
			result, err := linter.ParseAndLint(input)
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Issues()) != 0 {
				for _, issue := range result.Issues() {
					t.Error(issue.RuleName(), issue.Description())
				}
			}
			if result.Input() != input {
				t.Fatal("input changed")
			}
		})
	}
}

func TestLeadingBlankRulesRejectMissingCRLFBlank(t *testing.T) {
	_, validBody := (&rule.BodyLeadingBlankRule{}).Validate(&mockCommit{
		message: "feat: expose records\r\nBody text.", body: "Body text.",
	})
	_, validFooter := (&rule.FooterLeadingBlankRule{}).Validate(&mockCommit{
		message: "feat: expose records\r\nRefs #9", footer: "Refs #9",
	})
	if validBody || validFooter {
		t.Fatal("missing blank line accepted")
	}
}

func TestSyntaxFooterRules(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"footer-enum", "footer-type-enum", "signed-off-by", "trailer-exists"}
	conf.Settings["footer-enum"] = lint.RuleSetting{
		Argument: []string{"Signed-off-by", "Refs", "Trace"},
	}
	conf.Settings["footer-type-enum"] = lint.RuleSetting{
		Argument: []interface{}{
			map[interface{}]interface{}{
				"token":  "Refs",
				"types":  []string{"feat"},
				"values": []string{"#"},
			},
		},
	}
	conf.Settings["trailer-exists"] = lint.RuleSetting{Argument: "Trace"}
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatal(err)
	}

	passMsg := "feat: expose records\n\nTrace: ok\nSigned-off-by: Sam\nRefs: #9"
	res, err := linter.ParseAndLint(passMsg)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues()) != 0 {
		t.Fatalf("unexpected issues: %v", res.Issues())
	}

	failMsg := "feat: expose records\n\nUnknown: x\nRefs: bad"
	res, err = linter.ParseAndLint(failMsg)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues()) != 4 {
		t.Fatalf("got %d issues, want 4: %v", len(res.Issues()), res.Issues())
	}
}

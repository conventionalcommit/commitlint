package test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/conventionalcommit/commitlint/commit"
	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/lint/formatter"
)

func TestSyntaxDiagnosticsAndInput(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"header-trim"}
	conf.Severity.Default = lint.SeverityWarn
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name, input string
		wantRules   []string
	}{
		{"terminal newline", "feat: expose records\n", nil},
		{"crlf body", "feat: expose records\r\n\r\nBody text.\r\n", nil},
		{"trailing header spaces", "feat: expose records  \n", []string{"header-trim"}},
		{"leading space", " feat: expose records", []string{"parser", "header-trim"}},
		{"leading blank", "\nfeat: expose records", []string{"parser"}},
		{"missing section separator", "feat: expose records\nBody text.", []string{"parser"}},
		{"empty footer value", "feat: expose records\n\nRefs #9\nEmpty: ", []string{"parser"}},
		{"blank description", "feat: \t", []string{"parser", "header-trim"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := linter.ParseAndLint(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if result.Input() != tt.input {
				t.Fatal("result changed input bytes")
			}
			if len(result.Issues()) != len(tt.wantRules) {
				t.Fatalf("issues=%v, want %v", result.Issues(), tt.wantRules)
			}
			for i, wantRule := range tt.wantRules {
				issue := result.Issues()[i]
				severity := lint.SeverityWarn
				if wantRule == "parser" {
					severity = lint.SeverityError
				}
				if issue.RuleName() != wantRule || issue.Severity() != severity {
					t.Fatalf("issue=%s %s", issue.RuleName(), issue.Severity())
				}
			}
			text, err := (&formatter.JSONFormatter{}).Format(result)
			if err != nil {
				t.Fatal(err)
			}
			var output struct {
				Input string `json:"input"`
			}
			if err := json.Unmarshal([]byte(text), &output); err != nil {
				t.Fatal(err)
			}
			if output.Input != tt.input {
				t.Fatal("formatter changed input bytes")
			}
		})
	}
}

func TestSyntaxIgnoreCRLFAndDiagnosticLineColumn(t *testing.T) {
	linter := newDefaultLinter(t)
	for _, msg := range []string{"Initial commit\r\n", "Merge branch 'main'\r\nBody\r\n"} {
		res, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Issues()) != 0 {
			t.Fatalf("CRLF ignored message %q reported issues: %v", msg, res.Issues())
		}
	}

	res, err := linter.ParseAndLint("feat: expose records\n\nRefs #9\nEmpty: ")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues()) != 1 {
		t.Fatalf("issues = %v, want 1 parser issue", res.Issues())
	}
	infos := res.Issues()[0].Infos()
	if len(infos) != 1 || infos[0] != "line 4, column 8" {
		t.Fatalf("infos = %v, want [line 4, column 8]", infos)
	}
}

type customParser struct{}

func (p *customParser) Parse(msg string) (commit.Commit, error) {
	if msg == "bad" {
		return nil, errors.New("custom parse failure")
	}
	header, rest, _ := strings.Cut(msg, "\n\n")
	typ, desc, _ := strings.Cut(header, ": ")
	c := &mockCommit{
		message:     msg,
		header:      header,
		body:        rest,
		typ:         typ,
		description: desc,
	}
	if strings.HasPrefix(msg, "warn:") {
		return c, errors.New("recoverable custom error")
	}
	return c, nil
}

func TestCustomParserIntegration(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = append(conf.Rules, "body-leading-blank", "breaking-change-exclamation-mark", "signed-off-by")
	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := lint.NewWithParser(conf, rules, nil); err == nil {
		t.Fatal("expected error for nil parser in lint.NewWithParser")
	}
	if _, err := config.NewLinterWithParser(conf, nil); err == nil {
		t.Fatal("expected error for nil parser in config.NewLinterWithParser")
	}

	linter, err := config.NewLinterWithParser(conf, &customParser{})
	if err != nil {
		t.Fatal(err)
	}

	res, err := linter.ParseAndLint("bad")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues()) != 1 || res.Issues()[0].RuleName() != "parser" || len(res.Issues()[0].Infos()) != 0 {
		t.Fatalf("unexpected custom error result: %v", res.Issues())
	}

	res, err = linter.ParseAndLint("warn: partial commit")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues()) < 2 || res.Issues()[0].RuleName() != "parser" {
		t.Fatalf("expected parser error plus rule issues for partial commit, got: %v", res.Issues())
	}

	res, err = linter.ParseAndLint("feat: valid via custom parser\n\nBody line.")
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Issues()) != 1 || res.Issues()[0].RuleName() != "signed-off-by" {
		t.Fatalf("expected signed-off-by issue only, got: %v", res.Issues())
	}
}

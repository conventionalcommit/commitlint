package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

func TestSyntaxLineLengthBoundaries(t *testing.T) {
	for _, tt := range []struct {
		name, rule, prefix string
		limit              int
	}{
		{"body", "body-max-line-length", "", 72},
		{"footer", "footer-max-line-length", "Trace: ", 100},
	} {
		conf := config.NewDefault()
		conf.Rules = []string{tt.rule}
		conf.Settings[tt.rule] = lint.RuleSetting{Argument: tt.limit}
		linter, err := config.NewLinter(conf)
		if err != nil {
			t.Fatal(err)
		}
		for _, newline := range []string{"\n", "\r\n"} {
			for _, extra := range []int{0, 1} {
				t.Run(fmt.Sprintf("%s/%q/%d", tt.name, newline, extra), func(t *testing.T) {
					line := tt.prefix + strings.Repeat("x", tt.limit+extra-len(tt.prefix))
					next := "continuation"
					if tt.prefix != "" {
						next = "Refs #9"
					}
					input := "feat: expose records" + newline + newline + line + newline + next
					result, err := linter.ParseAndLint(input)
					if err != nil {
						t.Fatal(err)
					}
					if len(result.Issues()) != extra {
						t.Fatalf("issues=%v, want %d", result.Issues(), extra)
					}
					if extra == 1 {
						issue := result.Issues()[0]
						want := fmt.Sprintf("in line 1, length is %d", tt.limit+1)
						if issue.RuleName() != tt.rule || len(issue.Infos()) != 1 || issue.Infos()[0] != want {
							t.Fatalf("unexpected length issue: %s %v", issue.RuleName(), issue.Infos())
						}
					}
				})
			}
		}
	}
}

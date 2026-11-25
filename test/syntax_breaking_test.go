package test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/config"
)

func TestSyntaxBreakingIndicators(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"breaking-change-exclamation-mark"}
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name, input string
		wantIssue   bool
	}{
		{"neither", "feat: expose records", false},
		{"header only", "feat!: expose records", true},
		{"spaced footer only", "feat: expose records\n\nBREAKING CHANGE: remove key", true},
		{"hyphenated footer only", "feat: expose records\n\nBREAKING-CHANGE: remove key", true},
		{"both spaced", "feat!: expose records\n\nBREAKING CHANGE: remove key", false},
		{"both hyphenated", "feat(store)!: expose records\n\nBREAKING-CHANGE: remove key", false},
		{"description bang", "feat: retain ! in description", false},
		{"description bang with footer", "feat: retain ! in description\n\nBREAKING CHANGE: remove key", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := linter.ParseAndLint(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			want := 0
			if tt.wantIssue {
				want = 1
			}
			if len(result.Issues()) != want {
				t.Fatalf("issues=%v, want %d", result.Issues(), want)
			}
			if want != 0 && result.Issues()[0].RuleName() != "breaking-change-exclamation-mark" {
				t.Fatalf("unexpected issue: %s", result.Issues()[0].RuleName())
			}
		})
	}
}

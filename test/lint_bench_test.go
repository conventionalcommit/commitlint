package test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

var lintResult *lint.Result

func BenchmarkParseAndLint(b *testing.B) {
	linter, err := config.NewLinter(config.NewDefault())
	if err != nil {
		b.Fatal(err)
	}
	for _, tt := range []struct {
		name, input string
		invalid     bool
	}{
		{"header", "feat: expose records", false},
		{"body", "feat: expose records\n\nRetain records.", false},
		{"footers", "feat: expose records\n\nRetain records.\n\nRefs #9", false},
		{"recovery", " unknown:reload\nBody text.", true},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				lintResult, err = linter.ParseAndLint(tt.input)
				if err != nil {
					b.Fatal(err)
				}
			}
			if (len(lintResult.Issues()) > 0) != tt.invalid {
				b.Fatal("incorrect lint result")
			}
		})
	}
}

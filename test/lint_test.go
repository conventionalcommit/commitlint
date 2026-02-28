package test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

// ---------------------------------------------------------------------------
// Table-driven: valid commits (0 issues expected)
// ---------------------------------------------------------------------------

var validCommitCases = []struct {
	name string
	msg  string
}{
	// --- All default type-enum values ---
	{"type/feat", "feat: add new login feature"},
	{"type/fix", "fix: resolve crash on startup"},
	{"type/docs", "docs: update README with examples"},
	{"type/style", "style: format code with gofmt"},
	{"type/refactor", "refactor: simplify user service"},
	{"type/perf", "perf: optimize database queries"},
	{"type/test", "test: add unit tests for parser"},
	{"type/build", "build: update go version to 1.21"},
	{"type/ci", "ci: add GitHub Actions workflow"},
	{"type/chore", "chore: update all dependencies"},
	{"type/revert", "revert: undo last commit changes"},

	// --- Scope ---
	{"scope/simple", "feat(auth): add OAuth2 support"},
	{"scope/nested", "fix(auth/login): handle nil user"},
	{"scope/multi-word", "refactor(core): extract helper"},

	// --- Breaking indicator via exclamation mark ---
	{"breaking/excl-no-scope", "feat!: new breaking change"},
	{"breaking/excl-with-scope", "feat(api)!: change response format"},

	// --- Body ---
	{"body/single-line", "feat(auth): add OAuth2 support\n\nThis adds Google OAuth2."},
	{"body/multi-line", "fix: repair cache\n\nLine one.\nLine two.\nLine three."},

	// --- Footers (tokens the parser supports) ---
	{"footer/single", "feat: add feature\n\nSome body.\n\nFixes: #123"},
	{"footer/multi", "fix: something\n\nBody.\n\nFixes: #123\nReviewed-by: John"},
	{"footer/no-body", "feat: add utils\n\nFixes: #123"},

	// --- Header length boundaries ---
	// min=10 chars header; "ci: update" = 10 chars exactly
	{"header/at-min-10", "ci: update"},
	// max=72; "feat: " (6) + 66 = 72 chars exactly
	{"header/at-max-72", "feat: add a new feature that is at the exact seventy two character limit"},
}

func TestLint_ValidCommits(t *testing.T) {
	linter := newDefaultLinter(t)
	for _, tc := range validCommitCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linter.ParseAndLint(tc.msg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result.Issues()) != 0 {
				t.Errorf("expected 0 issues for %q, got %d:", tc.msg, len(result.Issues()))
				for _, iss := range result.Issues() {
					t.Logf("  [%s] %s: %s", iss.Severity(), iss.RuleName(), iss.Description())
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Table-driven: invalid commits (>= wantMin issues expected)
// ---------------------------------------------------------------------------

var invalidCommitCases = []struct {
	name    string
	msg     string
	wantMin int // minimum number of issues expected
}{
	// --- Parser errors ---
	{"parse/empty-message", "", 1},
	{"parse/no-colon-separator", "feat add new feature", 1},
	{"parse/missing-type", ": add new feature", 1},
	{"parse/space-in-type", "fe at: something here", 1},
	{"parse/missing-desc-no-space", "feat:", 1},
	{"parse/no-space-after-colon", "feat:compact description text", 1},
	{"parse/random-text", "just some random commit text", 1},

	// --- type-enum violations ---
	{"type-enum/unknown-type", "invalid: not an allowed type", 1},
	{"type-enum/cased-type", "Feat: case-sensitive type check", 1},
	{"type-enum/numeric-type", "123: something numeric", 1},

	// --- header-min-length violations (min=10) ---
	{"header-min/too-short", "fix: x", 1},
	// "fix: abcd" = 9 chars, under min of 10
	{"header-min/nine-chars", "fix: abcd", 1},

	// --- header-max-length violations (max=72) ---
	{"header-max/too-long", "feat: this is an extremely long commit message that exceeds the max header length limit", 1},
	// 73 chars: "feat: " (6) + 67 = 73
	{"header-max/one-over", "feat: add a new feature that is at the exact seventy two character limits", 1},
	// --- body-max-line-length violations (max=72) ---
	{"body-line/too-long", "feat: add feature\n\nThis is a very long body line that definitely exceeds the seventy-two character limit per line in body text", 1},
}

func TestLint_InvalidCommits(t *testing.T) {
	linter := newDefaultLinter(t)
	for _, tc := range invalidCommitCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linter.ParseAndLint(tc.msg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result.Issues()) < tc.wantMin {
				t.Errorf("expected >= %d issues for %q, got %d", tc.wantMin, tc.msg, len(result.Issues()))
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Specific rule-violation checks (verify the exact rule name in issues)
// ---------------------------------------------------------------------------

var ruleViolationCases = []struct {
	name     string
	msg      string
	wantRule string // expected RuleName in at least one issue
}{
	{"rule/type-enum", "invalid: not an allowed type", "type-enum"},
	{"rule/header-min-length", "fix: tiny", "header-min-length"},
	{"rule/header-max-length", "feat: this is an extremely long commit message that exceeds the max header length limit", "header-max-length"},
	{"rule/body-max-line-length", "feat: add feature\n\nThis is a very long body line that definitely exceeds the seventy-two character limit per line in body text", "body-max-line-length"},
	{"rule/footer-max-line-length", "feat: add utils\n\nBody text.\n\nFixes: this-is-a-very-long-footer-line-that-exceeds-the-one-hundred-character-limit-and-should-trigger", "footer-max-line-length"},
	{"rule/parser-error", "invalid message", "parser"},
}

func TestLint_RuleViolations(t *testing.T) {
	linter := newDefaultLinter(t)
	for _, tc := range ruleViolationCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := linter.ParseAndLint(tc.msg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			found := false
			for _, iss := range result.Issues() {
				if iss.RuleName() == tc.wantRule {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected rule %q in issues for %q", tc.wantRule, tc.msg)
				for _, iss := range result.Issues() {
					t.Logf("  got: [%s] %s: %s", iss.Severity(), iss.RuleName(), iss.Description())
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Result preservation
// ---------------------------------------------------------------------------

func TestLint_ResultPreservesInput(t *testing.T) {
	linter := newDefaultLinter(t)
	msgs := []string{
		"feat: add something new",
		"invalid message",
		"fix: x",
	}
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Input() != msg {
			t.Errorf("expected input %q, got %q", msg, result.Input())
		}
	}
}

// ---------------------------------------------------------------------------
// Severity tests
// ---------------------------------------------------------------------------

func TestLint_DefaultSeverityIsError(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("unknown: some message here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) == 0 {
		t.Fatal("expected at least one issue")
	}
	for _, iss := range result.Issues() {
		if iss.Severity() != lint.SeverityError {
			t.Errorf("expected error severity, got %q for %s", iss.Severity(), iss.RuleName())
		}
	}
}

func TestLint_CustomWarningSeverity(t *testing.T) {
	conf := config.NewDefault().Lint
	conf.Severity.Rules = map[string]lint.Severity{
		"type-enum": lint.SeverityWarn,
	}
	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatal(err)
	}
	linter, err := lint.New(conf, rules)
	if err != nil {
		t.Fatal(err)
	}

	result, err := linter.ParseAndLint("unknown: some message here")
	if err != nil {
		t.Fatal(err)
	}
	for _, iss := range result.Issues() {
		if iss.RuleName() == "type-enum" && iss.Severity() != lint.SeverityWarn {
			t.Errorf("expected warn for type-enum, got %q", iss.Severity())
		}
	}
}

func TestLint_ParserErrorAlwaysError(t *testing.T) {
	conf := config.NewDefault().Lint
	// Even with all rules set to warn, parser errors should be SeverityError
	conf.Severity.Default = lint.SeverityWarn
	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatal(err)
	}
	linter, err := lint.New(conf, rules)
	if err != nil {
		t.Fatal(err)
	}

	result, err := linter.ParseAndLint("")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Issues()) == 0 {
		t.Fatal("expected parser issue for empty msg")
	}
	for _, iss := range result.Issues() {
		if iss.RuleName() == "parser" && iss.Severity() != lint.SeverityError {
			t.Errorf("parser errors should be SeverityError, got %q", iss.Severity())
		}
	}
}

// ---------------------------------------------------------------------------
// Body edge cases
// ---------------------------------------------------------------------------

func TestLint_BodyWithinLimit(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("feat: add feature\n\nShort body line.")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Issues()) != 0 {
		t.Error("body within 72 chars should pass")
	}
}

func TestLint_BodyMultiLineAllWithin(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("feat: add feature\n\nLine one.\nLine two.\nLine three.")
	if err != nil {
		t.Fatal(err)
	}
	for _, iss := range result.Issues() {
		if iss.RuleName() == "body-max-line-length" {
			t.Error("all body lines <= 72 should not trigger body-max-line-length")
		}
	}
}

func TestLint_EmptyBodyAllowed(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("feat: add feature")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Issues()) != 0 {
		t.Error("empty body should be allowed by default config")
	}
}

// ---------------------------------------------------------------------------
// Multiple issues in a single commit
// ---------------------------------------------------------------------------

func TestLint_MultipleIssues(t *testing.T) {
	// "invalid: x" -> type-enum (invalid) + header-min-length (10 chars, "invalid: x" = 10)
	// Actually "invalid: x" = 10 chars so header-min passes. Use something shorter.
	// "xx: y" -> type: "xx" not in enum (5 chars header < 10)
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("xx: y")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Issues()) < 2 {
		t.Errorf("expected >= 2 issues (type-enum + header-min-length), got %d", len(result.Issues()))
		for _, iss := range result.Issues() {
			t.Logf("  %s: %s", iss.RuleName(), iss.Description())
		}
	}
}

package test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

// --- Ignore pattern tests ---

func TestIgnore_MergePullRequest(t *testing.T) {
	msgs := []string{
		"Merge pull request #123 from owner/branch",
		"Merge pull request #1 from a/b",
	}
	linter := newDefaultLinter(t)
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("error for %q: %v", msg, err)
		}
		if len(result.Issues()) != 0 {
			t.Errorf("expected %q to be ignored, got %d issues", msg, len(result.Issues()))
		}
	}
}

func TestIgnore_MergeBranch(t *testing.T) {
	msgs := []string{
		"Merge branch 'feature-x'",
		"Merge branch 'release/1.0' into main",
		"Merge feature into main",
		"Merge tag 'v1.0.0'",
		"Merge remote-tracking branch 'origin/main'",
		"Merge remote-tracking branch 'origin/main' into develop",
	}
	linter := newDefaultLinter(t)
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("error for %q: %v", msg, err)
		}
		if len(result.Issues()) != 0 {
			t.Errorf("expected %q to be ignored, got %d issues", msg, len(result.Issues()))
		}
	}
}

func TestIgnore_MergedPR(t *testing.T) {
	msgs := []string{
		"Merged PR #456: Add new feature",
		"Merged PR 789: Fix bug",
		"Merged feature-branch in main",
		"Merged feature-branch into main",
	}
	linter := newDefaultLinter(t)
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("error for %q: %v", msg, err)
		}
		if len(result.Issues()) != 0 {
			t.Errorf("expected %q to be ignored, got %d issues", msg, len(result.Issues()))
		}
	}
}

func TestIgnore_RevertReapply(t *testing.T) {
	msgs := []string{
		`Revert "feat: add new feature"`,
		`revert "fix: something"`,
		"Revert some commit",
		`Reapply "feat: add feature"`,
		`reapply "fix: something"`,
	}
	linter := newDefaultLinter(t)
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("error for %q: %v", msg, err)
		}
		if len(result.Issues()) != 0 {
			t.Errorf("expected %q to be ignored, got %d issues", msg, len(result.Issues()))
		}
	}
}

func TestIgnore_FixupSquashAmend(t *testing.T) {
	msgs := []string{
		"fixup! feat: add something",
		"squash! fix: repair something",
		"amend! chore: update deps",
	}
	linter := newDefaultLinter(t)
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("error for %q: %v", msg, err)
		}
		if len(result.Issues()) != 0 {
			t.Errorf("expected %q to be ignored, got %d issues", msg, len(result.Issues()))
		}
	}
}

func TestIgnore_AutomaticMerge(t *testing.T) {
	msgs := []string{
		"Automatic merge from release/1.0 to main",
		"Auto-merged feature-x into main",
	}
	linter := newDefaultLinter(t)
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("error for %q: %v", msg, err)
		}
		if len(result.Issues()) != 0 {
			t.Errorf("expected %q to be ignored, got %d issues", msg, len(result.Issues()))
		}
	}
}

func TestIgnore_InitialCommit(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("Initial commit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) != 0 {
		t.Errorf("expected 'Initial commit' to be ignored, got %d issues", len(result.Issues()))
	}
}

func TestIgnore_MultilineFirstLineMatches(t *testing.T) {
	msg := "Merge pull request #100 from org/branch\n\nAdditional body text"
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) != 0 {
		t.Errorf("expected multiline merge to be ignored, got %d issues", len(result.Issues()))
	}
}

func TestIgnore_NotIgnored(t *testing.T) {
	msgs := []string{
		"merge something random",
		"fixup something",
		"initial commit",
		"just some random text",
	}
	linter := newDefaultLinter(t)
	for _, msg := range msgs {
		result, err := linter.ParseAndLint(msg)
		if err != nil {
			t.Fatalf("error for %q: %v", msg, err)
		}
		// These should NOT be ignored, should produce issues
		if len(result.Issues()) == 0 {
			t.Errorf("expected %q to NOT be ignored (should have issues)", msg)
		}
	}
}

func TestIgnore_EmptyPatterns_NoSkip(t *testing.T) {
	conf := config.NewDefault().Lint
	conf.DisableDefaultIgnores = true
	conf.IgnorePatterns = []string{}

	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	linter, err := lint.New(conf, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := linter.ParseAndLint("Merge pull request #123 from owner/branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) == 0 {
		t.Error("expected merge to fail when all ignore patterns disabled")
	}
}

func TestIgnore_CustomPatternsAdditive(t *testing.T) {
	conf := config.NewDefault().Lint
	conf.IgnorePatterns = []string{`^CUSTOM-\d+`, `^WIP `}

	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	linter, err := lint.New(conf, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Custom pattern matches
	result, err := linter.ParseAndLint("CUSTOM-123 some ticket work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) != 0 {
		t.Errorf("expected custom pattern to be ignored, got %d issues", len(result.Issues()))
	}

	// WIP matches
	result, err = linter.ParseAndLint("WIP adding new feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) != 0 {
		t.Errorf("expected WIP to be ignored, got %d issues", len(result.Issues()))
	}

	// Default merge pattern STILL present (additive), should be ignored
	result, err = linter.ParseAndLint("Merge pull request #123 from owner/branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) != 0 {
		t.Error("merge should be ignored - user patterns are additive to defaults")
	}
}

func TestIgnore_CustomPatternsOnlyWhenDefaultsDisabled(t *testing.T) {
	conf := config.NewDefault().Lint
	conf.DisableDefaultIgnores = true
	conf.IgnorePatterns = []string{`^CUSTOM-\d+`}

	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	linter, err := lint.New(conf, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Custom pattern still matches
	result, err := linter.ParseAndLint("CUSTOM-123 some ticket work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) != 0 {
		t.Errorf("expected custom pattern to be ignored, got %d issues", len(result.Issues()))
	}

	// Default merge pattern should NOT be ignored
	result, err = linter.ParseAndLint("Merge pull request #123 from owner/branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues()) == 0 {
		t.Error("merge should fail when default ignores are disabled")
	}
}

func TestIgnore_InvalidPattern_LinterCreationFails(t *testing.T) {
	conf := config.NewDefault().Lint
	conf.IgnorePatterns = []string{`^valid`, `[invalid`}

	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("unexpected error getting rules: %v", err)
	}

	_, err = lint.New(conf, rules)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestIgnore_ValidationCatchesInvalidPattern(t *testing.T) {
	conf := config.NewDefault().Lint
	conf.IgnorePatterns = []string{`[invalid`}

	errs := config.ValidateLint(conf)
	found := false
	for _, e := range errs {
		if e != nil {
			found = true
		}
	}
	if !found {
		t.Error("expected validation error for invalid regex")
	}
}

func TestIgnore_DefaultPatternsExist(t *testing.T) {
	patterns := config.DefaultIgnorePatterns()
	if len(patterns) == 0 {
		t.Fatal("expected default ignore patterns to be non-empty")
	}
}

func TestIgnore_DefaultConfigHasPatterns(t *testing.T) {
	conf := config.NewDefault()
	if len(conf.Lint.DefaultIgnorePatterns) == 0 {
		t.Fatal("expected default config to have default ignore patterns")
	}
}

func TestIgnore_EffectiveIgnorePatterns(t *testing.T) {
	conf := config.NewDefault().Lint

	// Default: no user patterns, defaults enabled
	effective := conf.EffectiveIgnorePatterns()
	if len(effective) != len(conf.DefaultIgnorePatterns) {
		t.Errorf("expected %d effective patterns, got %d", len(conf.DefaultIgnorePatterns), len(effective))
	}

	// Add user patterns: effective = defaults + user
	conf.IgnorePatterns = []string{`^WIP `}
	effective = conf.EffectiveIgnorePatterns()
	expected := len(conf.DefaultIgnorePatterns) + 1
	if len(effective) != expected {
		t.Errorf("expected %d effective patterns, got %d", expected, len(effective))
	}

	// Disable defaults: effective = user only
	conf.DisableDefaultIgnores = true
	effective = conf.EffectiveIgnorePatterns()
	if len(effective) != 1 {
		t.Errorf("expected 1 effective pattern (defaults disabled), got %d", len(effective))
	}
}

func TestIgnore_IgnoredResultHasNoIssues(t *testing.T) {
	linter := newDefaultLinter(t)
	msg := "Merge pull request #42 from org/feature"
	result, err := linter.ParseAndLint(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Input() != msg {
		t.Errorf("expected input preserved, got %q", result.Input())
	}
	if len(result.Issues()) != 0 {
		t.Errorf("expected no issues for ignored commit, got %d", len(result.Issues()))
	}
}

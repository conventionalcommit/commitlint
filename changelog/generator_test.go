package changelog

import (
	"testing"

	"github.com/conventionalcommit/commitlint/internal/git"
)

func TestExtractReferencesIntegration(t *testing.T) {
	conf := &Config{
		IssuePrefixes: []string{"#", "GH-"},
	}

	msg := "fix: resolve #42 and GH-100"
	refs := git.ExtractReferences(msg, conf.IssuePrefixes)

	if len(refs) != 2 {
		t.Fatalf("expected 2 refs, got %d: %v", len(refs), refs)
	}
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"feat", "Feat"},
		{"", ""},
		{"Fix", "Fix"},
		{"ci", "Ci"},
	}

	for _, tt := range tests {
		result := capitalizeFirst(tt.input)
		if result != tt.expected {
			t.Errorf("capitalizeFirst(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestHasCommits(t *testing.T) {
	// nil
	if hasCommits(nil) {
		t.Error("hasCommits(nil) should be false")
	}

	// empty
	vc := &VersionChangelog{}
	if hasCommits(vc) {
		t.Error("hasCommits(empty) should be false")
	}

	// with groups
	vc.Groups = []TypeGroup{
		{Commits: []CommitInfo{{Hash: "abc"}}},
	}
	if !hasCommits(vc) {
		t.Error("hasCommits(with groups) should be true")
	}

	// with breaking
	vc2 := &VersionChangelog{
		Breaking: []CommitInfo{{Hash: "def"}},
	}
	if !hasCommits(vc2) {
		t.Error("hasCommits(with breaking) should be true")
	}
}

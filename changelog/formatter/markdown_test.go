package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/conventionalcommit/commitlint/changelog"
	"github.com/conventionalcommit/commitlint/commit"
)

func parseCommit(t *testing.T, msg string) commit.Commit {
	t.Helper()
	p := commit.NewParser()
	c, err := p.Parse(msg)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", msg, err)
	}
	return c
}

func TestMarkdownFormatterName(t *testing.T) {
	f := &MarkdownFormatter{}
	if f.Name() != "markdown" {
		t.Errorf("expected name 'markdown', got '%s'", f.Name())
	}
}

func TestMarkdownFormatVersion(t *testing.T) {
	f := &MarkdownFormatter{}
	v := &changelog.VersionChangelog{
		Version:    "v1.0.0",
		Date:       "2025-01-15",
		CompareURL: "https://github.com/u/r/compare/v0.9.0...v1.0.0",
		Groups: []changelog.TypeGroup{
			{
				Type:   "feat",
				Header: "Features",
				Commits: []changelog.CommitInfo{
					{
						Commit:    parseCommit(t, "feat: add something"),
						Hash:      "abc123def456",
						ShortHash: "abc123d",
						CommitURL: "https://github.com/u/r/commit/abc123def456",
						Date:      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	result, err := f.FormatVersion(v, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "## [v1.0.0]") {
		t.Error("expected version header with link")
	}
	if !strings.Contains(result, "(2025-01-15)") {
		t.Error("expected date in header")
	}
	if !strings.Contains(result, "### Features") {
		t.Error("expected Features header")
	}
	if !strings.Contains(result, "add something") {
		t.Error("expected commit description")
	}
	if !strings.Contains(result, "[abc123d]") {
		t.Error("expected short hash link")
	}
}

func TestMarkdownFormatVersionWithScope(t *testing.T) {
	f := &MarkdownFormatter{}
	v := &changelog.VersionChangelog{
		Version: "v1.0.0",
		Date:    "2025-01-15",
		Groups: []changelog.TypeGroup{
			{
				Type:   "feat",
				Header: "Features",
				Commits: []changelog.CommitInfo{
					{
						Commit:    parseCommit(t, "feat(api): add endpoint"),
						ShortHash: "abc1234",
						CommitURL: "https://example.com/commit/abc",
					},
				},
			},
		},
	}

	result, err := f.FormatVersion(v, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "**api:**") {
		t.Error("expected scope in bold")
	}
}

func TestMarkdownFormatBreakingChanges(t *testing.T) {
	f := &MarkdownFormatter{}
	v := &changelog.VersionChangelog{
		Version: "v2.0.0",
		Date:    "2025-06-01",
		Breaking: []changelog.CommitInfo{
			{
				Commit:    parseCommit(t, "feat!: breaking api change"),
				ShortHash: "def456",
			},
		},
		Groups: []changelog.TypeGroup{
			{
				Type:   "feat",
				Header: "Features",
				Commits: []changelog.CommitInfo{
					{
						Commit:    parseCommit(t, "feat!: breaking api change"),
						ShortHash: "def456",
					},
				},
			},
		},
	}

	result, err := f.FormatVersion(v, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(result, "⚠ BREAKING CHANGES") {
		t.Error("expected breaking changes section")
	}
}

func TestMarkdownFullFormat(t *testing.T) {
	f := &MarkdownFormatter{}
	cl := &changelog.Changelog{
		Header: "# Changelog",
		Versions: []changelog.VersionChangelog{
			{
				Version: "v1.0.0",
				Date:    "2025-01-15",
				Groups: []changelog.TypeGroup{
					{
						Type:   "feat",
						Header: "Features",
						Commits: []changelog.CommitInfo{
							{
								Commit:    parseCommit(t, "feat: first feature"),
								ShortHash: "aaa1111",
							},
						},
					},
				},
			},
		},
	}

	result, err := f.Format(cl)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(result, "# Changelog") {
		t.Error("expected changelog header at start")
	}
	if !strings.Contains(result, "## v1.0.0") {
		t.Error("expected version header")
	}
}

func TestMarkdownNoCommitURL(t *testing.T) {
	f := &MarkdownFormatter{}
	v := &changelog.VersionChangelog{
		Version: "v1.0.0",
		Date:    "2025-01-15",
		Groups: []changelog.TypeGroup{
			{
				Type:   "feat",
				Header: "Features",
				Commits: []changelog.CommitInfo{
					{
						Commit:    parseCommit(t, "feat: something"),
						ShortHash: "abc1234",
						// no CommitURL
					},
				},
			},
		},
	}

	result, err := f.FormatVersion(v, nil)
	if err != nil {
		t.Fatal(err)
	}

	// should have (abc1234) without link
	if !strings.Contains(result, "(abc1234)") {
		t.Error("expected short hash without link")
	}
}

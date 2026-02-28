package formatter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/conventionalcommit/commitlint/changelog"
)

func TestJSONFormatterName(t *testing.T) {
	f := &JSONFormatter{}
	if f.Name() != "json" {
		t.Errorf("expected name 'json', got '%s'", f.Name())
	}
}

func TestJSONFormatVersion(t *testing.T) {
	f := &JSONFormatter{}

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
						Author:    "John",
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

	// verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed["version"] != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got %v", parsed["version"])
	}
}

func TestJSONFullFormat(t *testing.T) {
	f := &JSONFormatter{}

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
								Author:    "Jane",
								Date:      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
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

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if parsed["header"] != "# Changelog" {
		t.Error("expected header in JSON")
	}

	versions, ok := parsed["versions"].([]interface{})
	if !ok || len(versions) != 1 {
		t.Error("expected 1 version in JSON")
	}
}

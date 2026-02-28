package git

import (
	"testing"
)

func TestNormalizeRemoteURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "github ssh",
			input:    "git@github.com:user/repo.git",
			expected: "https://github.com/user/repo",
		},
		{
			name:     "github https with .git",
			input:    "https://github.com/user/repo.git",
			expected: "https://github.com/user/repo",
		},
		{
			name:     "github https without .git",
			input:    "https://github.com/user/repo",
			expected: "https://github.com/user/repo",
		},
		{
			name:     "gitlab ssh",
			input:    "git@gitlab.com:group/project.git",
			expected: "https://gitlab.com/group/project",
		},
		{
			name:     "bitbucket ssh",
			input:    "git@bitbucket.org:team/repo.git",
			expected: "https://bitbucket.org/team/repo",
		},
		{
			name:     "http to https",
			input:    "http://github.com/user/repo.git",
			expected: "https://github.com/user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeRemoteURL(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeRemoteURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetectHostType(t *testing.T) {
	tests := []struct {
		url      string
		expected HostType
	}{
		{"https://github.com/user/repo", HostGitHub},
		{"https://gitlab.com/group/project", HostGitLab},
		{"https://bitbucket.org/team/repo", HostBitbucket},
		{"https://dev.azure.com/org/project", HostAzure},
		{"https://visualstudio.com/org/project", HostAzure},
		{"https://custom.example.com/repo", HostGeneric},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			result := DetectHostType(tt.url)
			if result != tt.expected {
				t.Errorf("detectHostType(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestInferCommitURL(t *testing.T) {
	base := "https://github.com/user/repo"
	result := InferCommitURL(base, HostGitHub)
	expected := "https://github.com/user/repo/commit/{{hash}}"
	if result != expected {
		t.Errorf("inferCommitURL() = %q, want %q", result, expected)
	}

	result = InferCommitURL("https://gitlab.com/g/p", HostGitLab)
	expected = "https://gitlab.com/g/p/-/commit/{{hash}}"
	if result != expected {
		t.Errorf("inferCommitURL() = %q, want %q", result, expected)
	}
}

func TestInferCompareURL(t *testing.T) {
	result := InferCompareURL("https://github.com/u/r", HostGitHub)
	expected := "https://github.com/u/r/compare/{{from}}...{{to}}"
	if result != expected {
		t.Errorf("inferCompareURL() = %q, want %q", result, expected)
	}

	result = InferCompareURL("https://dev.azure.com/o/p", HostAzure)
	expected = "https://dev.azure.com/o/p#/compare?head=true&sourceBranch={{to}}&targetBranch={{from}}"
	if result != expected {
		t.Errorf("inferCompareURL() = %q, want %q", result, expected)
	}
}

func TestBuildCommitURL(t *testing.T) {
	tmpl := "https://github.com/user/repo/commit/{{hash}}"
	result := BuildCommitURL(tmpl, "abc123")
	expected := "https://github.com/user/repo/commit/abc123"
	if result != expected {
		t.Errorf("buildCommitURL() = %q, want %q", result, expected)
	}

	// empty template
	result = BuildCommitURL("", "abc123")
	if result != "" {
		t.Errorf("buildCommitURL empty template = %q, want empty", result)
	}
}

func TestBuildCompareURL(t *testing.T) {
	tmpl := "https://github.com/user/repo/compare/{{from}}...{{to}}"
	result := BuildCompareURL(tmpl, "v1.0.0", "v2.0.0")
	expected := "https://github.com/user/repo/compare/v1.0.0...v2.0.0"
	if result != expected {
		t.Errorf("buildCompareURL() = %q, want %q", result, expected)
	}
}

func TestParseGitLog(t *testing.T) {
	output := commitSep + "\nabc123full" + fieldSep + "abc1234" + fieldSep + "John Doe" + fieldSep + "2025-01-15T10:00:00+00:00" + fieldSep + "feat: add feature" + fieldSep + "some body"

	commits := parseGitLog(output)
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}

	c := commits[0]
	if c.Hash != "abc123full" {
		t.Errorf("hash = %q, want %q", c.Hash, "abc123full")
	}
	if c.ShortHash != "abc1234" {
		t.Errorf("shortHash = %q, want %q", c.ShortHash, "abc1234")
	}
	if c.Author != "John Doe" {
		t.Errorf("author = %q, want %q", c.Author, "John Doe")
	}
	if c.Subject != "feat: add feature" {
		t.Errorf("subject = %q, want %q", c.Subject, "feat: add feature")
	}
	if c.Body != "some body" {
		t.Errorf("body = %q, want %q", c.Body, "some body")
	}
}

func TestExtractReferences(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		prefixes []string
		expected []string
	}{
		{
			name:     "single reference",
			msg:      "fix: resolve issue #42",
			prefixes: []string{"#"},
			expected: []string{"#42"},
		},
		{
			name:     "multiple references",
			msg:      "fix: resolve #42 and #100",
			prefixes: []string{"#"},
			expected: []string{"#42", "#100"},
		},
		{
			name:     "no references",
			msg:      "fix: resolve issue",
			prefixes: []string{"#"},
			expected: nil,
		},
		{
			name:     "custom prefix",
			msg:      "fix: resolve JIRA-123",
			prefixes: []string{"JIRA-"},
			expected: []string{"JIRA-123"},
		},
		{
			name:     "no duplicates",
			msg:      "fix: resolve #42 and #42",
			prefixes: []string{"#"},
			expected: []string{"#42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refs := ExtractReferences(tt.msg, tt.prefixes)
			if len(refs) != len(tt.expected) {
				t.Errorf("extractReferences() returned %d refs, want %d", len(refs), len(tt.expected))
				return
			}
			for i, ref := range refs {
				if ref != tt.expected[i] {
					t.Errorf("ref[%d] = %q, want %q", i, ref, tt.expected[i])
				}
			}
		})
	}
}

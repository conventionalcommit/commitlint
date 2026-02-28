package test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/conventionalcommit/commitlint/lint/formatter"
)

func TestDefaultFormatter_Name(t *testing.T) {
	f := &formatter.DefaultFormatter{}
	if f.Name() != "default" {
		t.Errorf("expected 'default', got %q", f.Name())
	}
}

func TestDefaultFormatter_NoIssues(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("feat: valid commit message")
	if err != nil {
		t.Fatal(err)
	}
	f := &formatter.DefaultFormatter{}
	out, fErr := f.Format(result)
	if fErr != nil {
		t.Fatal(fErr)
	}
	if !strings.Contains(out, "\u2714") {
		t.Errorf("expected checkmark for valid commit, got %q", out)
	}
}

func TestDefaultFormatter_WithIssues(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("invalid commit no type")
	if err != nil {
		t.Fatal(err)
	}
	f := &formatter.DefaultFormatter{}
	out, fErr := f.Format(result)
	if fErr != nil {
		t.Fatal(fErr)
	}
	if strings.Contains(out, "\u2714") {
		t.Errorf("should not have checkmark for invalid commit, got %q", out)
	}
	if !strings.Contains(out, "commitlint") {
		t.Errorf("expected 'commitlint' header in output, got %q", out)
	}
}

func TestJSONFormatter_Name(t *testing.T) {
	f := &formatter.JSONFormatter{}
	if f.Name() != "json" {
		t.Errorf("expected 'json', got %q", f.Name())
	}
}

func TestJSONFormatter_NoIssues(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("feat: valid commit message")
	if err != nil {
		t.Fatal(err)
	}
	f := &formatter.JSONFormatter{}
	out, fErr := f.Format(result)
	if fErr != nil {
		t.Fatal(fErr)
	}
	var parsed map[string]interface{}
	if jErr := json.Unmarshal([]byte(out), &parsed); jErr != nil {
		t.Fatalf("invalid JSON output: %v", jErr)
	}
	issues, ok := parsed["issues"].([]interface{})
	if !ok {
		t.Fatal("expected 'issues' array in JSON output")
	}
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %d", len(issues))
	}
}

func TestJSONFormatter_WithIssues(t *testing.T) {
	linter := newDefaultLinter(t)
	result, err := linter.ParseAndLint("invalid commit no type")
	if err != nil {
		t.Fatal(err)
	}
	f := &formatter.JSONFormatter{}
	out, fErr := f.Format(result)
	if fErr != nil {
		t.Fatal(fErr)
	}
	var parsed map[string]interface{}
	if jErr := json.Unmarshal([]byte(out), &parsed); jErr != nil {
		t.Fatalf("invalid JSON output: %v", jErr)
	}
	issues, ok := parsed["issues"].([]interface{})
	if !ok {
		t.Fatal("expected 'issues' array in JSON")
	}
	if len(issues) == 0 {
		t.Error("expected at least one issue for invalid commit")
	}
	// Verify each issue has mandatory fields
	for i, raw := range issues {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("issue[%d] is not a map", i)
		}
		for _, field := range []string{"name", "severity", "description"} {
			if _, exists := entry[field]; !exists {
				t.Errorf("issue[%d] missing field %q", i, field)
			}
		}
	}
}

func TestJSONFormatter_InputField(t *testing.T) {
	linter := newDefaultLinter(t)
	msg := "feat: check input field"
	result, err := linter.ParseAndLint(msg)
	if err != nil {
		t.Fatal(err)
	}
	f := &formatter.JSONFormatter{}
	out, fErr := f.Format(result)
	if fErr != nil {
		t.Fatal(fErr)
	}
	var parsed map[string]interface{}
	if jErr := json.Unmarshal([]byte(out), &parsed); jErr != nil {
		t.Fatal(jErr)
	}
	if parsed["input"] != msg {
		t.Errorf("expected input %q, got %q", msg, parsed["input"])
	}
}

func TestDefaultFormatter_TruncatesLongInput(t *testing.T) {
	linter := newDefaultLinter(t)
	// Create a commit that will fail (no type) but has a very long first line
	longMsg := "this is a very long commit message without a type that will definitely be truncated in the output"
	result, err := linter.ParseAndLint(longMsg)
	if err != nil {
		t.Fatal(err)
	}
	f := &formatter.DefaultFormatter{}
	out, fErr := f.Format(result)
	if fErr != nil {
		t.Fatal(fErr)
	}
	// The default formatter truncates input to 25 chars
	if !strings.Contains(out, "...") {
		t.Errorf("expected truncation with '...' in output, got %q", out)
	}
}

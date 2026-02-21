package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

func TestConfig_NewDefault(t *testing.T) {
	conf := config.NewDefault()

	if conf.MinVersion == "" {
		t.Error("expected non-empty MinVersion")
	}
	if conf.Formatter == "" {
		t.Error("expected non-empty Formatter")
	}
	if conf.Formatter != "default" {
		t.Errorf("expected formatter 'default', got %q", conf.Formatter)
	}
	if len(conf.Rules) == 0 {
		t.Error("expected non-empty rules")
	}
	if len(conf.Settings) == 0 {
		t.Error("expected non-empty settings")
	}
	if conf.Severity.Default != lint.SeverityError {
		t.Errorf("expected default severity 'error', got %q", conf.Severity.Default)
	}
	if len(conf.DefaultIgnorePatterns) == 0 {
		t.Error("expected non-empty default ignore patterns")
	}
	if len(conf.IgnorePatterns) != 0 {
		t.Error("expected empty user ignore patterns in default config")
	}
}

func TestConfig_NewLinter(t *testing.T) {
	conf := config.NewDefault()
	linter, err := config.NewLinter(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if linter == nil {
		t.Fatal("expected non-nil linter")
	}
}

func TestConfig_GetFormatter(t *testing.T) {
	conf := config.NewDefault()
	f, err := config.GetFormatter(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
	if f.Name() != "default" {
		t.Errorf("expected formatter name 'default', got %q", f.Name())
	}
}

func TestConfig_GetFormatterJSON(t *testing.T) {
	conf := config.NewDefault()
	conf.Formatter = "json"
	f, err := config.GetFormatter(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Name() != "json" {
		t.Errorf("expected formatter name 'json', got %q", f.Name())
	}
}

func TestConfig_GetFormatterUnknown(t *testing.T) {
	conf := config.NewDefault()
	conf.Formatter = "unknown"
	_, err := config.GetFormatter(conf)
	if err == nil {
		t.Error("expected error for unknown formatter")
	}
}

func TestConfig_Validate_Valid(t *testing.T) {
	conf := config.NewDefault()
	errs := config.Validate(conf)
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got %d:", len(errs))
		for _, e := range errs {
			t.Errorf("  - %v", e)
		}
	}
}

func TestConfig_Validate_InvalidFormatter(t *testing.T) {
	conf := config.NewDefault()
	conf.Formatter = "nonexistent"
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation errors for unknown formatter")
	}
}

func TestConfig_Validate_EmptyFormatter(t *testing.T) {
	conf := config.NewDefault()
	conf.Formatter = ""
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation errors for empty formatter")
	}
}

func TestConfig_Validate_InvalidSeverity(t *testing.T) {
	conf := config.NewDefault()
	conf.Severity.Default = "invalid"
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation errors for invalid severity")
	}
}

func TestConfig_Validate_InvalidRuleSeverity(t *testing.T) {
	conf := config.NewDefault()
	conf.Severity.Rules = map[string]lint.Severity{
		"type-enum": "invalid-severity",
	}
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation errors for invalid rule severity")
	}
}

func TestConfig_Validate_UnknownRule(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = append(conf.Rules, "nonexistent-rule")
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation errors for unknown rule")
	}
}

func TestConfig_Validate_InvalidIgnorePattern(t *testing.T) {
	conf := config.NewDefault()
	conf.IgnorePatterns = []string{`[invalid`}
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation errors for invalid ignore pattern")
	}
}

func TestConfig_Validate_ValidIgnorePattern(t *testing.T) {
	conf := config.NewDefault()
	conf.IgnorePatterns = []string{`^Merge .*`}
	errs := config.Validate(conf)
	for _, e := range errs {
		t.Errorf("unexpected validation error: %v", e)
	}
}

func TestConfig_WriteTo(t *testing.T) {
	conf := config.NewDefault()
	var buf bytes.Buffer
	err := config.WriteCompactTo(&buf, conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output from WriteTo")
	}

	// Should only contain settings for the 5 enabled rules, not all 20
	for _, disabled := range []string{"body-min-length", "body-max-length", "type-min-length", "scope-charset"} {
		if bytes.Contains(buf.Bytes(), []byte(disabled+":")) {
			t.Errorf("should not contain disabled rule setting %q", disabled)
		}
	}
	for _, enabled := range []string{"header-min-length", "header-max-length", "type-enum"} {
		if !bytes.Contains(buf.Bytes(), []byte(enabled+":")) {
			t.Errorf("should contain enabled rule setting %q", enabled)
		}
	}
}

func TestConfig_WriteTo_WithUserIgnores(t *testing.T) {
	conf := config.NewDefault()
	conf.IgnorePatterns = []string{`^WIP `}
	var buf bytes.Buffer
	err := config.WriteCompactTo(&buf, conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("ignores")) {
		t.Error("expected 'ignores' field in YAML output when user patterns exist")
	}
}

func TestConfig_Parse_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	confContent := `min-version: v0.9.0
formatter: default
rules:
  - header-min-length
  - header-max-length
  - type-enum
severity:
  default: error
settings:
  header-min-length:
    argument: 10
  header-max-length:
    argument: 50
  type-enum:
    argument:
      - feat
      - fix
`
	err := os.WriteFile(confPath, []byte(confContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	conf, err := config.Parse(confPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if conf.Formatter != "default" {
		t.Errorf("expected formatter 'default', got %q", conf.Formatter)
	}
	if len(conf.Rules) != 3 {
		t.Errorf("expected 3 rules, got %d", len(conf.Rules))
	}
}

func TestConfig_Parse_OldVersionKey(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	// Uses the old "version:" key for backward compatibility
	confContent := `version: v0.9.0
formatter: default
rules:
  - header-min-length
severity:
  default: error
settings:
  header-min-length:
    argument: 10
`
	err := os.WriteFile(confPath, []byte(confContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	conf, err := config.Parse(confPath)
	if err != nil {
		t.Fatalf("unexpected error parsing old 'version' key: %v", err)
	}

	if conf.MinVersion != "v0.9.0" {
		t.Errorf("expected MinVersion 'v0.9.0', got %q", conf.MinVersion)
	}
}

func TestConfig_Parse_WithIgnores(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	confContent := `min-version: v0.9.0
formatter: default
rules:
  - header-min-length
severity:
  default: error
settings:
  header-min-length:
    argument: 10
ignores:
  - "^WIP "
  - "^TICKET-\\d+"
`
	err := os.WriteFile(confPath, []byte(confContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	conf, err := config.Parse(confPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conf.IgnorePatterns) != 2 {
		t.Errorf("expected 2 ignore patterns, got %d", len(conf.IgnorePatterns))
	}
}

func TestConfig_Parse_WithoutIgnores_UsesDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	confContent := `min-version: v0.9.0
formatter: default
rules:
  - header-min-length
severity:
  default: error
settings:
  header-min-length:
    argument: 10
`
	err := os.WriteFile(confPath, []byte(confContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	conf, err := config.Parse(confPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(conf.DefaultIgnorePatterns) == 0 {
		t.Error("expected default ignore patterns to be populated by Parse")
	}
	if len(conf.IgnorePatterns) != 0 {
		t.Error("expected empty user ignore patterns when not specified in config")
	}
}

func TestConfig_Parse_DisableDefaultIgnores(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	confContent := `min-version: v0.9.0
formatter: default
rules:
  - header-min-length
severity:
  default: error
settings:
  header-min-length:
    argument: 10
disable-default-ignores: true
ignores:
  - "^WIP "
`
	err := os.WriteFile(confPath, []byte(confContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	conf, err := config.Parse(confPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !conf.DisableDefaultIgnores {
		t.Error("expected DisableDefaultIgnores to be true")
	}
	if len(conf.IgnorePatterns) != 1 {
		t.Errorf("expected 1 user ignore pattern, got %d", len(conf.IgnorePatterns))
	}
	effective := conf.EffectiveIgnorePatterns()
	if len(effective) != 1 {
		t.Errorf("expected 1 effective pattern (defaults disabled), got %d", len(effective))
	}
}

func TestConfig_Parse_NonExistentFile(t *testing.T) {
	_, err := config.Parse("/nonexistent/path/commitlint.yaml")
	if err == nil {
		t.Error("expected error for non-existent config file")
	}
}

func TestConfig_Parse_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	err := os.WriteFile(confPath, []byte("invalid: yaml: content: ["), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err = config.Parse(confPath)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestConfig_GetEnabledRules(t *testing.T) {
	conf := config.NewDefault()
	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != len(conf.Rules) {
		t.Errorf("expected %d rules, got %d", len(conf.Rules), len(rules))
	}
}

func TestConfig_GetEnabledRules_DuplicateRules(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = append(conf.Rules, conf.Rules[0])

	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != len(conf.Rules)-1 {
		t.Errorf("expected duplicates to be removed, got %d rules", len(rules))
	}
}

func TestConfig_GetEnabledRules_UnknownRule(t *testing.T) {
	conf := config.NewDefault()
	conf.Rules = []string{"nonexistent-rule"}

	_, err := config.GetEnabledRules(conf)
	if err == nil {
		t.Error("expected error for unknown rule")
	}
}

func TestConfig_SeverityString(t *testing.T) {
	tests := []struct {
		severity lint.Severity
		expected string
	}{
		{lint.SeverityError, "Error"},
		{lint.SeverityWarn, "Warning"},
		{lint.Severity("unknown"), "Severity(unknown)"},
	}

	for _, tc := range tests {
		got := tc.severity.String()
		if got != tc.expected {
			t.Errorf("Severity(%q).String() = %q, want %q", tc.severity, got, tc.expected)
		}
	}
}

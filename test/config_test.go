package test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/conventionalcommit/commitlint/changelog"
	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

func TestConfig_NewDefault(t *testing.T) {
	conf := config.NewDefault().Lint

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
	linter, err := config.NewLinter(conf.Lint)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if linter == nil {
		t.Fatal("expected non-nil linter")
	}
}

func TestConfig_GetFormatter(t *testing.T) {
	conf := config.NewDefault()
	f, err := config.GetFormatter(conf.Lint)
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
	conf.Lint.Formatter = "json"
	f, err := config.GetFormatter(conf.Lint)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Name() != "json" {
		t.Errorf("expected formatter name 'json', got %q", f.Name())
	}
}

func TestConfig_GetFormatterUnknown(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.Formatter = "unknown"
	_, err := config.GetFormatter(conf.Lint)
	if err == nil {
		t.Error("expected error for unknown formatter")
	}
}

func TestConfig_Validate_Valid(t *testing.T) {
	conf := config.NewDefault()
	errs := config.ValidateLint(conf.Lint)
	if len(errs) != 0 {
		t.Errorf("expected no validation errors, got %d:", len(errs))
		for _, e := range errs {
			t.Errorf("  - %v", e)
		}
	}
}

func TestConfig_Validate_InvalidFormatter(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.Formatter = "nonexistent"
	errs := config.ValidateLint(conf.Lint)
	if len(errs) == 0 {
		t.Error("expected validation errors for unknown formatter")
	}
}

func TestConfig_Validate_EmptyFormatter(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.Formatter = ""
	errs := config.ValidateLint(conf.Lint)
	if len(errs) == 0 {
		t.Error("expected validation errors for empty formatter")
	}
}

func TestConfig_Validate_InvalidSeverity(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.Severity.Default = "invalid"
	errs := config.ValidateLint(conf.Lint)
	if len(errs) == 0 {
		t.Error("expected validation errors for invalid severity")
	}
}

func TestConfig_Validate_InvalidRuleSeverity(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.Severity.Rules = map[string]lint.Severity{
		"type-enum": "invalid-severity",
	}
	errs := config.ValidateLint(conf.Lint)
	if len(errs) == 0 {
		t.Error("expected validation errors for invalid rule severity")
	}
}

func TestConfig_Validate_UnknownRule(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.Rules = append(conf.Lint.Rules, "nonexistent-rule")
	errs := config.ValidateLint(conf.Lint)
	if len(errs) == 0 {
		t.Error("expected validation errors for unknown rule")
	}
}

func TestConfig_Validate_InvalidIgnorePattern(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.IgnorePatterns = []string{`[invalid`}
	errs := config.ValidateLint(conf.Lint)
	if len(errs) == 0 {
		t.Error("expected validation errors for invalid ignore pattern")
	}
}

func TestConfig_Validate_ValidIgnorePattern(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.IgnorePatterns = []string{`^Merge .*`}
	errs := config.ValidateLint(conf.Lint)
	for _, e := range errs {
		t.Errorf("unexpected validation error: %v", e)
	}
}

func TestConfig_WriteCompactTo(t *testing.T) {
	conf := config.NewDefault()
	var buf bytes.Buffer
	err := config.WriteCompactTo(&buf, conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output from WriteCompactTo")
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

	// Should not contain hidden changelog types (e.g. chore, style are hidden by default)
	for _, hidden := range []string{"chore", "style", "refactor", "test", "build", "ci", "revert"} {
		if bytes.Contains(buf.Bytes(), []byte("type: "+hidden)) {
			t.Errorf("compact output should not contain hidden type %q", hidden)
		}
	}
	// Visible types should still be present
	for _, visible := range []string{"feat", "fix", "docs", "perf"} {
		if !bytes.Contains(buf.Bytes(), []byte("type: "+visible)) {
			t.Errorf("compact output should contain visible type %q", visible)
		}
	}
}

func TestConfig_WriteCompactTo_WithUserIgnores(t *testing.T) {
	conf := config.NewDefault()
	conf.Lint.IgnorePatterns = []string{`^WIP `}
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

	confContent := `lint:
  min-version: v0.9.0
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

	if conf.Lint.Formatter != "default" {
		t.Errorf("expected formatter 'default', got %q", conf.Lint.Formatter)
	}
	if len(conf.Lint.Rules) != 3 {
		t.Errorf("expected 3 rules, got %d", len(conf.Lint.Rules))
	}
}

func TestConfig_Parse_OldVersionKey(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	// Uses the old "version:" key for backward compatibility
	confContent := `lint:
  version: v0.9.0
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

	if conf.Lint.MinVersion != "v0.9.0" {
		t.Errorf("expected MinVersion 'v0.9.0', got %q", conf.Lint.MinVersion)
	}
}

func TestConfig_Parse_WithIgnores(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	confContent := `lint:
  min-version: v0.9.0
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

	if len(conf.Lint.IgnorePatterns) != 2 {
		t.Errorf("expected 2 ignore patterns, got %d", len(conf.Lint.IgnorePatterns))
	}
}

func TestConfig_Parse_WithoutIgnores_UsesDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	confContent := `lint:
  min-version: v0.9.0
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

	if len(conf.Lint.DefaultIgnorePatterns) == 0 {
		t.Error("expected default ignore patterns to be populated by Parse")
	}
	if len(conf.Lint.IgnorePatterns) != 0 {
		t.Error("expected empty user ignore patterns when not specified in config")
	}
}

func TestConfig_Parse_DisableDefaultIgnores(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	confContent := `lint:
  min-version: v0.9.0
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

	if !conf.Lint.DisableDefaultIgnores {
		t.Error("expected DisableDefaultIgnores to be true")
	}
	if len(conf.Lint.IgnorePatterns) != 1 {
		t.Errorf("expected 1 user ignore pattern, got %d", len(conf.Lint.IgnorePatterns))
	}
	effective := conf.Lint.EffectiveIgnorePatterns()
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
	conf := config.NewDefault().Lint
	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != len(conf.Rules) {
		t.Errorf("expected %d rules, got %d", len(conf.Rules), len(rules))
	}
}

func TestConfig_GetEnabledRules_DuplicateRules(t *testing.T) {
	conf := config.NewDefault().Lint
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
	conf := config.NewDefault().Lint
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

// --- Old config format detection ---

func TestConfig_Parse_OldFlatConfig(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	// Old pre-v0.12.0 flat config (no "lint:" wrapper)
	confContent := `formatter: default
rules:
  - header-min-length
  - type-enum
severity:
  default: error
settings:
  header-min-length:
    argument: 10
  type-enum:
    argument:
      - feat
      - fix
`
	err := os.WriteFile(confPath, []byte(confContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err = config.Parse(confPath)
	if err == nil {
		t.Fatal("expected error for old flat config format")
	}
	if !strings.Contains(err.Error(), "pre-v0.12.0") {
		t.Errorf("expected migration hint in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "migration.md") {
		t.Errorf("expected migration.md link in error, got: %v", err)
	}
}

// --- Changelog defaults in Parse ---

func TestConfig_Parse_ChangelogDefaultsApplied(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	// Config with only lint section, no changelog
	confContent := `lint:
  min-version: v0.9.0
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

	if conf.Changelog == nil {
		t.Fatal("expected non-nil changelog config")
	}
	if conf.Changelog.Formatter != "markdown" {
		t.Errorf("expected default changelog formatter 'markdown', got %q", conf.Changelog.Formatter)
	}
	if conf.Changelog.Header != "# Changelog" {
		t.Errorf("expected default changelog header, got %q", conf.Changelog.Header)
	}
	if len(conf.Changelog.Types) == 0 {
		t.Error("expected default changelog types")
	}
	if len(conf.Changelog.IssuePrefixes) == 0 {
		t.Error("expected default issue prefixes")
	}
}

func TestConfig_Parse_PartialChangelog(t *testing.T) {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "commitlint.yaml")

	// Config with partial changelog (only formatter)
	confContent := `lint:
  min-version: v0.9.0
  formatter: default
  rules:
    - header-min-length
  severity:
    default: error
  settings:
    header-min-length:
      argument: 10
changelog:
  formatter: json
`
	err := os.WriteFile(confPath, []byte(confContent), 0o644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	conf, err := config.Parse(confPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Explicit value should be kept
	if conf.Changelog.Formatter != "json" {
		t.Errorf("expected changelog formatter 'json', got %q", conf.Changelog.Formatter)
	}
	// Missing fields should get defaults
	if conf.Changelog.Header != "# Changelog" {
		t.Errorf("expected default changelog header, got %q", conf.Changelog.Header)
	}
	if len(conf.Changelog.Types) == 0 {
		t.Error("expected default changelog types when not specified")
	}
}

// --- WriteTo / WriteCompactTo with nil ---

func TestConfig_WriteTo_NilChangelog(t *testing.T) {
	conf := &config.Config{
		Lint:      config.NewDefaultLint(),
		Changelog: nil,
	}
	var buf bytes.Buffer
	err := config.WriteTo(&buf, conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
	// Should contain changelog section from defaults
	if !bytes.Contains(buf.Bytes(), []byte("changelog:")) {
		t.Error("expected 'changelog:' section in output")
	}
}

func TestConfig_WriteCompactTo_NilChangelog(t *testing.T) {
	conf := &config.Config{
		Lint:      config.NewDefaultLint(),
		Changelog: nil,
	}
	var buf bytes.Buffer
	err := config.WriteCompactTo(&buf, conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestConfig_WriteTo_NilLint(t *testing.T) {
	conf := &config.Config{
		Lint:      nil,
		Changelog: config.NewDefaultChangelog(),
	}
	var buf bytes.Buffer
	err := config.WriteTo(&buf, conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("lint:")) {
		t.Error("expected 'lint:' section in output")
	}
}

// --- Validate with nil ---

func TestConfig_Validate_NilLint(t *testing.T) {
	conf := &config.Config{
		Lint:      nil,
		Changelog: config.NewDefaultChangelog(),
	}
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation error for nil lint config")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "lint config is nil") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'lint config is nil' error")
	}
}

func TestConfig_Validate_NilChangelog(t *testing.T) {
	conf := &config.Config{
		Lint:      config.NewDefaultLint(),
		Changelog: nil,
	}
	errs := config.Validate(conf)
	if len(errs) == 0 {
		t.Error("expected validation error for nil changelog config")
	}
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "changelog config is nil") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'changelog config is nil' error")
	}
}

func TestConfig_Validate_Full(t *testing.T) {
	conf := config.NewDefault()
	errs := config.Validate(conf)
	if len(errs) != 0 {
		t.Errorf("expected no validation errors for default config, got %d:", len(errs))
		for _, e := range errs {
			t.Errorf("  - %v", e)
		}
	}
}

// --- ValidateChangelog ---

func TestConfig_ValidateChangelog_Valid(t *testing.T) {
	conf := config.NewDefaultChangelog()
	errs := config.ValidateChangelog(conf)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %d", len(errs))
	}
}

func TestConfig_ValidateChangelog_EmptyFormatter(t *testing.T) {
	conf := config.NewDefaultChangelog()
	conf.Formatter = ""
	errs := config.ValidateChangelog(conf)
	if len(errs) == 0 {
		t.Error("expected error for empty formatter")
	}
}

func TestConfig_ValidateChangelog_UnknownFormatter(t *testing.T) {
	conf := config.NewDefaultChangelog()
	conf.Formatter = "nonexistent"
	errs := config.ValidateChangelog(conf)
	if len(errs) == 0 {
		t.Error("expected error for unknown formatter")
	}
}

func TestConfig_ValidateChangelog_EmptyTypes(t *testing.T) {
	conf := config.NewDefaultChangelog()
	conf.Types = nil
	errs := config.ValidateChangelog(conf)
	if len(errs) == 0 {
		t.Error("expected error for empty types")
	}
}

func TestConfig_ValidateChangelog_DuplicateType(t *testing.T) {
	conf := config.NewDefaultChangelog()
	conf.Types = append(conf.Types, conf.Types[0]) // duplicate first type
	errs := config.ValidateChangelog(conf)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "duplicate") {
			found = true
		}
	}
	if !found {
		t.Error("expected duplicate type error")
	}
}

func TestConfig_ValidateChangelog_EmptyTypeField(t *testing.T) {
	conf := config.NewDefaultChangelog()
	conf.Types = append(conf.Types, changelog.TypeConfig{Type: "", Header: "Empty"})
	errs := config.ValidateChangelog(conf)
	found := false
	for _, e := range errs {
		if strings.Contains(e.Error(), "empty type field") {
			found = true
		}
	}
	if !found {
		t.Error("expected 'empty type field' error")
	}
}

// --- GetChangelogFormatter ---

func TestConfig_GetChangelogFormatter(t *testing.T) {
	conf := config.NewDefaultChangelog()
	f, err := config.GetChangelogFormatter(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
	if f.Name() != "markdown" {
		t.Errorf("expected 'markdown', got %q", f.Name())
	}
}

func TestConfig_GetChangelogFormatter_Unknown(t *testing.T) {
	conf := config.NewDefaultChangelog()
	conf.Formatter = "nonexistent"
	_, err := config.GetChangelogFormatter(conf)
	if err == nil {
		t.Error("expected error for unknown formatter")
	}
}

func TestConfig_GetChangelogFormatter_Empty(t *testing.T) {
	conf := config.NewDefaultChangelog()
	conf.Formatter = ""
	_, err := config.GetChangelogFormatter(conf)
	if err == nil {
		t.Error("expected error for empty formatter")
	}
}

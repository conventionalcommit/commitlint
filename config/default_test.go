package config

import (
	"testing"

	"github.com/conventionalcommit/commitlint/registry"
)

func TestDefaultLint(t *testing.T) {
	defConf := NewDefault().Lint
	_, err := NewLinter(defConf)
	if err != nil {
		t.Error("default lint creation failed", err)
		return
	}
}

func TestDefaultSettings(t *testing.T) {
	defConf := NewDefault()
	rules := registry.Rules()
	settingSize := len(defConf.Lint.Settings)
	if len(rules) != settingSize {
		t.Error("default config does not have all rule settings", len(rules), settingSize)
		return
	}
}

func TestNewLintDefault(t *testing.T) {
	conf := NewDefaultLint()
	if conf.MinVersion == "" {
		t.Error("expected non-empty MinVersion")
	}
	if conf.Formatter == "" {
		t.Error("expected non-empty Formatter")
	}
	if len(conf.Rules) == 0 {
		t.Error("expected non-empty rules")
	}
	if len(conf.Settings) == 0 {
		t.Error("expected non-empty settings")
	}
}

func TestNewDefaultChangelog_Valid(t *testing.T) {
	conf := NewDefaultChangelog()
	errs := ValidateChangelog(conf)
	if len(errs) != 0 {
		t.Errorf("expected no validation errors for default changelog, got %d:", len(errs))
		for _, e := range errs {
			t.Errorf("  - %v", e)
		}
	}
}

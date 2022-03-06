package config

import (
	"testing"

	"github.com/conventionalcommit/commitlint/internal/registry"
)

func TestDefaultLint(t *testing.T) {
	defConf := NewDefault()
	_, err := NewLinter(defConf)
	if err != nil {
		t.Error("default lint creation failed", err)
		return
	}
}

func TestDefaultSettings(t *testing.T) {
	defConf := NewDefault()
	rules := registry.Rules()
	settingSize := len(defConf.Settings)
	if len(rules) != settingSize {
		t.Error("default config does not have all rule settings", len(rules), settingSize)
		return
	}
}

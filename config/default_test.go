package config

import (
	"testing"

	"github.com/conventionalcommit/commitlint/internal/registry"
)

func TestDefaultLint(t *testing.T) {
	_, err := NewLinter(defConf)
	if err != nil {
		t.Error("default lint creation failed", err)
		return
	}
}

func TestDefaultConf(t *testing.T) {
	rules := registry.Rules()
	if len(rules) != len(defConf.Rules) {
		t.Error("default conf does not have all rules", len(rules), len(defConf.Rules))
		return
	}
}

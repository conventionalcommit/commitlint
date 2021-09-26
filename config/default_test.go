package config

import (
	"testing"
)

func TestDefaultLint(t *testing.T) {
	_, err := GetLinter(defConf)
	if err != nil {
		t.Error("default lint creation failed", err)
		return
	}
}

func TestDefaultConf(t *testing.T) {
	if len(defaultRules) != len(defConf.Rules) {
		t.Error("default conf does not have all rules", len(defaultRules), len(defConf.Rules))
		return
	}
}

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
	rules := globalRegistry.Rules()
	if len(rules) != len(defConf.Rules) {
		t.Error("default conf does not have all rules", len(rules), len(defConf.Rules))
		return
	}
}

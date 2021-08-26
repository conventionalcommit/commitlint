package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLint(t *testing.T) {
	_, err := GetLinter(defConf)
	assert.NoError(t, err, "default lint creation failed")
}

func TestDefaultConf(t *testing.T) {
	if len(allRules) != len(defConf.Rules) {
		t.Error("default conf does not have all rules", len(allRules), len(defConf.Rules))
		return
	}
}

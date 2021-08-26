package config

import "testing"

func TestAllRules(t *testing.T) {
	var m = make(map[string]struct{})
	for _, r := range allRules {
		_, ok := m[r.Name()]
		if ok {
			t.Errorf("error: %s rule name already exists", r.Name())
		}
		m[r.Name()] = struct{}{}
	}
}

func TestAllFormatters(t *testing.T) {
	var m = make(map[string]struct{})
	for _, r := range allFormatters {
		_, ok := m[r.Name()]
		if ok {
			t.Errorf("error: %s formatter name already exists", r.Name())
		}
		m[r.Name()] = struct{}{}
	}
}

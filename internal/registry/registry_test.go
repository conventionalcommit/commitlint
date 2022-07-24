package registry

import "testing"

func TestDefaultRules(t *testing.T) {
	m := make(map[string]struct{})
	for _, r := range globalRegistry.Rules() {
		_, ok := m[r.Name()]
		if ok {
			t.Errorf("error: %s rule name already exists", r.Name())
		}
		m[r.Name()] = struct{}{}
	}
}

func TestDefaultFormatters(t *testing.T) {
	m := make(map[string]struct{})
	for _, r := range globalRegistry.Formatters() {
		_, ok := m[r.Name()]
		if ok {
			t.Errorf("error: %s formatter name already exists", r.Name())
		}
		m[r.Name()] = struct{}{}
	}
}

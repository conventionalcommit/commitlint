package registry

import "testing"

func TestDefaultRulesNoDuplicates(t *testing.T) {
	m := make(map[string]struct{})
	for _, r := range Rules() {
		if _, ok := m[r.Name()]; ok {
			t.Errorf("duplicate rule name: %s", r.Name())
		}
		m[r.Name()] = struct{}{}
	}
}

func TestDefaultFormattersNoDuplicates(t *testing.T) {
	m := make(map[string]struct{})
	for _, f := range Formatters() {
		if _, ok := m[f.Name()]; ok {
			t.Errorf("duplicate formatter name: %s", f.Name())
		}
		m[f.Name()] = struct{}{}
	}
}

func TestRegisterCustomRule(t *testing.T) {
	// Registering an already-registered rule must return an error.
	rules := Rules()
	if len(rules) == 0 {
		t.Fatal("expected at least one default rule")
	}
	err := RegisterRule(rules[0])
	if err == nil {
		t.Errorf("expected error when registering duplicate rule %q", rules[0].Name())
	}
}

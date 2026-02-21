package test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

// mockCommit implements lint.Commit for testing
type mockCommit struct {
	message     string
	header      string
	body        string
	footer      string
	typ         string
	scope       string
	description string
	notes       []lint.Note
	breaking    bool
}

func (m *mockCommit) Message() string        { return m.message }
func (m *mockCommit) Header() string         { return m.header }
func (m *mockCommit) Body() string           { return m.body }
func (m *mockCommit) Footer() string         { return m.footer }
func (m *mockCommit) Type() string           { return m.typ }
func (m *mockCommit) Scope() string          { return m.scope }
func (m *mockCommit) Description() string    { return m.description }
func (m *mockCommit) Notes() []lint.Note     { return m.notes }
func (m *mockCommit) IsBreakingChange() bool { return m.breaking }

// mockNote implements lint.Note for testing
type mockNote struct {
	token string
	value string
}

func (n *mockNote) Token() string { return n.token }
func (n *mockNote) Value() string { return n.value }

// newDefaultLinter creates a linter with default config for testing
func newDefaultLinter(t *testing.T) *lint.Linter {
	t.Helper()
	conf := config.NewDefault()
	rules, err := config.GetEnabledRules(conf)
	if err != nil {
		t.Fatalf("failed to get enabled rules: %v", err)
	}
	linter, err := lint.New(conf, rules)
	if err != nil {
		t.Fatalf("failed to create linter: %v", err)
	}
	return linter
}

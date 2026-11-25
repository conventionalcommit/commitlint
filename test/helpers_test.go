package test

import (
	"strings"
	"testing"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

// mockCommit implements lint.Commit for testing without depending on any parser package.
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

func (m *mockCommit) HasHeaderBang() bool {
	colon := strings.IndexByte(m.header, ':')
	return colon > 0 && m.header[colon-1] == '!'
}

func (m *mockCommit) HasBreakingFooter() bool {
	for _, n := range m.notes {
		if t := n.Token(); t == "BREAKING CHANGE" || t == "BREAKING-CHANGE" {
			return true
		}
	}
	return false
}

func (m *mockCommit) BodyBlankLines() int {
	if m.body == "" {
		return 0
	}
	idx := strings.IndexByte(m.message, '\n')
	if idx < 0 {
		return 0
	}
	rest := m.message[idx+1:]
	if strings.HasPrefix(rest, "\n") || strings.HasPrefix(rest, "\r\n") {
		return 1
	}
	return 0
}

func (m *mockCommit) FooterBlankLines() int {
	if m.footer == "" {
		return 0
	}
	idx := strings.LastIndex(m.message, m.footer)
	if idx <= 0 || m.message[idx-1] != '\n' {
		return 0
	}
	prev := idx - 1
	if prev > 0 && m.message[prev-1] == '\r' {
		prev--
	}
	if prev > 0 && m.message[prev-1] == '\n' {
		return 1
	}
	return 0
}

func (m *mockCommit) NoteCount() int         { return len(m.notes) }
func (m *mockCommit) NoteAt(i int) lint.Note { return m.notes[i] }

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

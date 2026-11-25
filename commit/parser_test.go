package commit_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/conventionalcommit/commitlint/commit"
	"github.com/conventionalcommit/syntax"
)

func TestSyntaxAdapter(t *testing.T) {
	p := commit.NewParser()
	for _, newline := range []string{"\n", "\r\n"} {
		input := strings.ReplaceAll("feat(store)!: replace index  \n\nRetain records.\n\nTrace: #91\n continuation\nRefs #92\nRefs #93\nBREAKING-CHANGE: remove legacy keys\n", "\n", newline)
		c, err := p.Parse(input)
		if err != nil {
			t.Fatal(err)
		}
		if c.Message() != input || c.Header() != "feat(store)!: replace index  " || c.Type() != "feat" || c.Scope() != "store" || c.Description() != "replace index  " || c.Body() != "Retain records." || !c.IsBreakingChange() {
			t.Fatalf("unexpected fields: %#v", c)
		}
		want := [][2]string{{"Trace", "#91" + newline + " continuation"}, {"Refs", "92"}, {"Refs", "93"}, {"BREAKING-CHANGE", "remove legacy keys"}}
		notes := c.Notes()
		if len(notes) != len(want) {
			t.Fatalf("notes = %d", len(notes))
		}
		for i, entry := range want {
			if notes[i].Token() != entry[0] || notes[i].Value() != entry[1] {
				t.Fatalf("note %d = %q: %q", i, notes[i].Token(), notes[i].Value())
			}
		}
		notes[0] = nil
		if c.Notes()[0].Token() != "Trace" {
			t.Fatal("returned slice changed the result")
		}
		if !strings.HasPrefix(c.Footer(), "Trace: #91"+newline) {
			t.Fatalf("footer = %q", c.Footer())
		}
	}
}

func TestSyntaxAdapterErrors(t *testing.T) {
	p := commit.NewParser()
	for _, tt := range []struct {
		input  string
		code   syntax.ErrorCode
		offset int
	}{
		{" feat: expose records", syntax.InvalidType, 0},
		{"feat:expose records", syntax.InvalidSeparator, 5},
		{"feat: expose records\nbody", syntax.InvalidSectionSeparator, 21},
		{"feat: expose records\n\nRefs #9\nEmpty: ", syntax.InvalidFooterValue, 37},
	} {
		c, err := p.Parse(tt.input)
		var pe *syntax.ParseError
		if c != nil || !errors.As(err, &pe) || pe.Code != tt.code || pe.Offset != tt.offset {
			t.Fatalf("input %q: result=%v error=%v", tt.input, c, err)
		}
	}
}

func TestSyntaxAdapterWithoutNotes(t *testing.T) {
	c, err := commit.NewParser().Parse("feat: expose records\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Notes()) != 0 || c.Message() != "feat: expose records\n" {
		t.Fatal("unexpected empty-footer result")
	}
}

func TestSyntaxInspectionAdapter(t *testing.T) {
	p := commit.NewParser()
	input := "feat:expose records\nBody text."
	c, diagnostics := commit.Inspect(input)
	if c == nil || len(diagnostics) != 2 || c.Type() != "feat" || c.Body() != "Body text." {
		t.Fatalf("result=%v diagnostics=%v", c, diagnostics)
	}
	var pe *syntax.ParseError
	if !errors.As(diagnostics[0], &pe) || pe.Code != syntax.InvalidSeparator {
		t.Fatal("typed diagnostic lost")
	}
	if strict, err := p.Parse(input); strict != nil || err == nil {
		t.Fatal("strict parsing recovered invalid input")
	}
	if c, diagnostics := commit.Inspect("invalid header"); c != nil || len(diagnostics) != 1 {
		t.Fatal("ambiguous header recovered")
	}
}

func TestCommitHelpers(t *testing.T) {
	parsed, err := commit.NewParser().Parse("feat!: drop v1\n\nBody text.\n\nBREAKING CHANGE: gone\nRefs #9")
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.HasHeaderBang() || !parsed.HasBreakingFooter() || parsed.BodyBlankLines() != 1 || parsed.FooterBlankLines() != 1 {
		t.Fatal("syntax-backed Commit methods reported incorrect facts on valid breaking commit")
	}
	if parsed.NoteCount() != 2 || parsed.NoteAt(0).Token() != "BREAKING CHANGE" || parsed.NoteAt(1).Token() != "Refs" {
		t.Fatalf("unexpected notes: count=%d", parsed.NoteCount())
	}
	if got := commit.FirstLine("Initial commit\r\nignored"); got != "Initial commit" {
		t.Fatalf("FirstLine = %q", got)
	}
	var lines []string
	commit.ForEachLine("a\r\nb\nc", func(line string) bool {
		lines = append(lines, line)
		return len(lines) < 2
	})
	if len(lines) != 2 || lines[0] != "a" || lines[1] != "b" {
		t.Fatalf("ForEachLine = %q", lines)
	}
	_, diags := commit.Inspect("feat: ok\n\nRefs #9\nEmpty: ")
	if line, col, ok := commit.DiagnosticPosition("feat: ok\n\nRefs #9\nEmpty: ", diags[0]); !ok || line != 4 || col != 8 {
		t.Fatalf("DiagnosticPosition = (%d, %d, %v)", line, col, ok)
	}
	if _, _, ok := commit.DiagnosticPosition("feat: ok", errors.New("custom error")); ok {
		t.Fatal("expected false for non-syntax error")
	}
}

var parsed commit.Commit
var parseErr error
var parsedNotes []commit.Note

func BenchmarkSyntaxAdapter(b *testing.B) {
	p := commit.NewParser()
	for _, tt := range []struct {
		name, input string
		notes       bool
	}{
		{"header", "feat: expose records", false},
		{"footers", "feat: expose records\n\nRefs #9\nTrace: #91", false},
		{"footers_with_notes", "feat: expose records\n\nRefs #9\nTrace: #91", true},
	} {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				parsed, parseErr = p.Parse(tt.input)
				if tt.notes && parseErr == nil {
					parsedNotes = parsed.Notes()
				}
			}
			if parseErr != nil {
				b.Fatal(parseErr)
			}
		})
	}
}

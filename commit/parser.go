package commit

import (
	"errors"

	"github.com/conventionalcommit/syntax"
)

type defaultParser struct{}

// NewParser returns a stateless parser backed by syntax.
func NewParser() Parser {
	return &defaultParser{}
}

func (p *defaultParser) Parse(input string) (Commit, error) {
	c, err := syntax.Parse(input)
	if err != nil {
		return nil, err
	}
	return &defaultCommit{Commit: c}, nil
}

// Inspect exposes recoverable fields alongside all syntax diagnostics.
// A nil commit means structural recovery was unsafe.
func Inspect(input string) (Commit, []error) {
	result := syntax.Inspect(input)
	var diagnostics []error
	if len(result.Diagnostics) > 0 {
		diagnostics = make([]error, len(result.Diagnostics))
		for i := range diagnostics {
			diagnostics[i] = &result.Diagnostics[i]
		}
	}
	if !result.Structured {
		return nil, diagnostics
	}
	return &defaultCommit{Commit: result.Commit}, diagnostics
}

type defaultCommit struct {
	syntax.Commit
}

func (d *defaultCommit) BodyBlankLines() int {
	return d.BodyGap().BlankLines
}

func (d *defaultCommit) FooterBlankLines() int {
	return d.FooterGap().BlankLines
}

func (d *defaultCommit) NoteCount() int {
	return d.FooterCount()
}

func (d *defaultCommit) NoteAt(index int) Note {
	entry := d.FooterAt(index)
	return defaultNote{
		token: entry.Token(),
		value: entry.Value(),
	}
}

func (d *defaultCommit) Notes() []Note {
	count := d.FooterCount()
	if count == 0 {
		return nil
	}
	// Keep one backing slice for all notes instead of allocating each note.
	values := make([]defaultNote, count)
	notes := make([]Note, count)
	for i := range values {
		entry := d.FooterAt(i)
		values[i] = defaultNote{
			token: entry.Token(),
			value: entry.Value(),
		}
		notes[i] = &values[i]
	}
	return notes
}

type defaultNote struct {
	token string
	value string
}

func (n defaultNote) Token() string { return n.token }
func (n defaultNote) Value() string { return n.value }

// FirstLine returns the first line of msg without its LF or CRLF ending.
func FirstLine(msg string) string {
	it := syntax.Lines(msg)
	line, _ := it.Next()
	return msg[line.Content.Start:line.Content.End]
}

// ForEachLine iterates over lines in s without their LF or CRLF endings.
// Iteration stops if fn returns false.
func ForEachLine(s string, fn func(line string) bool) {
	it := syntax.Lines(s)
	for line, ok := it.Next(); ok; line, ok = it.Next() {
		if !fn(s[line.Content.Start:line.Content.End]) {
			return
		}
	}
}

// DiagnosticPosition returns the 1-based line and column for err in msg when
// err carries a byte offset, or false otherwise.
func DiagnosticPosition(msg string, err error) (line, column int, ok bool) {
	var pe *syntax.ParseError
	if !errors.As(err, &pe) {
		return 0, 0, false
	}
	pos := syntax.Locate(msg, pe.Offset)
	return pos.Line, pos.Column, true
}

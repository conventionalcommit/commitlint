// Package commit provides core interfaces for parsed conventional commits.
// Consumers depend on these interfaces independently of the backing parser.
package commit

// Note represents a footer note.
type Note interface {
	Token() string
	Value() string
}

// Commit represents a parsed conventional commit message.
type Commit interface {
	Message() string
	Header() string
	Body() string
	Footer() string
	Type() string
	Scope() string
	Description() string
	Notes() []Note
	IsBreakingChange() bool

	HasHeaderBang() bool
	HasBreakingFooter() bool
	BodyBlankLines() int
	FooterBlankLines() int
	NoteCount() int
	NoteAt(index int) Note
}

// Parser parses a commit message into a Commit.
type Parser interface {
	Parse(msg string) (Commit, error)
}

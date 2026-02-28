// Package commit provides core interfaces for parsed conventional commits.
// It wraps the external parser package and is the single import point for it.
// Both the lint and changelog packages depend on these interfaces.
package commit

// Note represent a footer note
type Note interface {
	Token() string
	Value() string
}

// Commit represent a parsed conventional commit message
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
}

// Parser parses a commit message into a Commit
type Parser interface {
	Parse(msg string) (Commit, error)
}

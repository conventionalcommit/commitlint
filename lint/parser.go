package lint

import "github.com/conventionalcommit/parser"

// Commit represent a commit message
// for now it is an alias of parser.Commit
type Commit = parser.Commit

// Parser parses given commit message
type Parser interface {
	Parse(msg string) (*Commit, error)
}

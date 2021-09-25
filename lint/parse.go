package lint

import (
	"github.com/conventionalcommit/parser"
)

// Parse parses given commit message
func Parse(msg string) (*Commit, error) {
	return parser.Parse(msg)
}

func isHeaderErr(err error) bool {
	return parser.IsHeaderErr(err)
}

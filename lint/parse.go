package lint

import (
	"github.com/conventionalcommit/parser"
)

// Parse parses given msg and checks for header error
func Parse(msg string) (*Commit, bool, error) {
	commit, err := parser.Parse(msg)
	if err != nil {
		return nil, parser.IsHeaderErr(err), err
	}
	return commit, false, nil
}

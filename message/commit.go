// Package message contains commit message
package message

import "github.com/conventionalcommit/parser"

// Commit is alias for parser.Commit
type Commit = parser.Commit

// Parse parses given msg and checks for header error
func Parse(msg string) (*Commit, bool, error) {
	commit, err := parser.Parse(msg)
	if err != nil {
		return nil, parser.IsHeaderErr(err), err
	}
	return commit, false, nil
}

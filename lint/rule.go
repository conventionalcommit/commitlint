package lint

import "github.com/conventionalcommit/commitlint/message"

// Rule represent a linter rule
type Rule interface {
	// Name returns name of the rule, it should be a unique identifier
	Name() string

	// Apply sets the argument to the rule from config file
	// if args are invalid or not expected return an error
	// Apply is called before Validate
	Apply(arg interface{}) error

	// Validate validates the rule for given message
	Validate(msg *message.Commit) (result string, isValid bool)
}

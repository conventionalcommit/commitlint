package lint

// Rule represent a linter rule
type Rule interface {
	// Name returns name of the rule, it should be a unique identifier
	Name() string

	// Apply calls with arguments and flags for the rule from config file
	// if flags or arguments are invalid or not expected return an error
	// Apply is called before Validate
	Apply(arg interface{}, flags map[string]interface{}) error

	// Validate validates the rule for given commit message
	// if given commit is valid, return true and messages slice are ignored
	// if invalid, return a error messages with false
	Validate(msg *Commit) (messages []string, isValid bool)
}

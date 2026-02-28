package lint

import "github.com/conventionalcommit/commitlint/commit"

// Type aliases for backward compatibility.
// External packages that depend on lint.Commit, lint.Note, or lint.Parser
// continue to work without changes.
type (
	Commit = commit.Commit
	Note   = commit.Note
	Parser = commit.Parser
)

// Rule Severity Constants
const (
	SeverityWarn  Severity = "warn"
	SeverityError Severity = "error"
)

// Severity represent the severity level of a rule
type Severity string

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "Error"
	case SeverityWarn:
		return "Warning"
	default:
		return "Severity(" + string(s) + ")"
	}
}

// Formatter represent a lint result formatter
type Formatter interface {
	// Name is a unique identifier for formatter
	Name() string

	// Format formats the linter result
	Format(result *Result) (string, error)
}

// Rule represent a linter rule
type Rule interface {
	// Name returns name of the rule, it should be a unique identifier
	Name() string

	// Apply calls with arguments and flags for the rule from config file
	// if flags or arguments are invalid or not expected return an error
	// Apply is called before Validate
	Apply(setting RuleSetting) error

	// Validate validates the rule for given commit message
	// if given commit is valid, return true and messages slice are ignored
	// if invalid, return a error messages with false
	Validate(msg Commit) (issue *Issue, isValid bool)
}

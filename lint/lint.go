package lint

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

// Note represent a footer note
type Note interface {
	Token() string
	Value() string
}

// Commit represent a commit message
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

// Parser parses given commit message
type Parser interface {
	Parse(msg string) (Commit, error)
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
	Validate(msg Commit) (messages []string, isValid bool)
}

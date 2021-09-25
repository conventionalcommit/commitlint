package lint

// Formatter represent a lint result formatter
type Formatter interface {
	// Name is a unique identifier for formatter
	Name() string

	// Format formats the linter result
	Format(result *Result) (string, error)
}

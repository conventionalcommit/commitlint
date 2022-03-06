package lint

// Result holds a linter result
type Result struct {
	input  string
	issues []*Issue
}

func newResult(input string, issues ...*Issue) *Result {
	return &Result{
		input:  input,
		issues: issues,
	}
}

// Input returns the input commit message
func (r *Result) Input() string { return r.input }

// Issues returns linter issues
func (r *Result) Issues() []*Issue { return r.issues }

// Issue holds a rule result
type Issue struct {
	ruleName string

	severity Severity

	description string

	additionalInfos []string
}

// NewIssue returns a new issue
func NewIssue(desc string, infos ...string) *Issue {
	return &Issue{
		description:     desc,
		additionalInfos: infos,
	}
}

// RuleName returns rule name
func (r *Issue) RuleName() string { return r.ruleName }

// Severity returns severity of the Rule Failure
func (r *Issue) Severity() Severity { return r.severity }

// Description returns description of the issue
func (r *Issue) Description() string { return r.description }

// Infos returns additional infos about the issue
func (r *Issue) Infos() []string { return r.additionalInfos }

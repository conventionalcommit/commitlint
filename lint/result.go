package lint

// Result holds Result of linter
type Result struct {
	inputMsg string

	issues []*Issue
}

func newResult(inputMsg string) *Result {
	return &Result{
		inputMsg: inputMsg,
	}
}

// AddError adds
func (res *Result) add(r *Issue) {
	res.issues = append(res.issues, r)
}

// IsOK returns true if commit message passed all the rules
func (res *Result) IsOK() bool { return len(res.issues) == 0 }

// Input returns input commit message
func (res *Result) Input() string { return res.inputMsg }

// Issues returns rule Issues
func (res *Result) Issues() []*Issue { return res.issues }

// Issue holds Failure of a linter rule
type Issue struct {
	name     string
	severity Severity
	messages []string
}

func newIssue(name string, msgs []string, severity Severity) *Issue {
	return &Issue{
		name:     name,
		messages: msgs,
		severity: severity,
	}
}

// Name returns rule name
func (r *Issue) Name() string { return r.name }

// Severity returns severity of the Rule Failure
func (r *Issue) Severity() Severity { return r.severity }

// Message returns the error messages of failed rule
func (r *Issue) Message() []string { return r.messages }

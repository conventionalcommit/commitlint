package lint

// Result holds result of linter
type Result struct {
	inputMsg string

	errResults  []RuleResult
	warnResults []RuleResult
}

// RuleResult holds result of a linter rule
type RuleResult struct {
	Name     string
	Severity string
	Message  string
}

func newResult(inputMsg string) *Result {
	return &Result{inputMsg: inputMsg}
}

// AddError adds
func (res *Result) add(r RuleResult) {
	if r.Severity == SeverityError {
		res.errResults = append(res.errResults, r)
	} else if r.Severity == SeverityWarn {
		res.warnResults = append(res.warnResults, r)
	} else {
		// should not come here
		res.errResults = append(res.errResults, r)
	}
}

// Input returns input commit message
func (res *Result) Input() string { return res.inputMsg }

// Errors returns all error messages
func (res *Result) Errors() []RuleResult { return res.errResults }

// Warns returns all warning messages
func (res *Result) Warns() []RuleResult { return res.warnResults }

// HasErrors returns true if errors found by linter
func (res *Result) HasErrors() bool { return len(res.errResults) != 0 }

// HasWarns returns true if warnings found by linter
func (res *Result) HasWarns() bool { return len(res.warnResults) != 0 }

// IsOK returns true if commit message passed all the rules
func (res *Result) IsOK() bool { return !res.HasErrors() && !res.HasWarns() }

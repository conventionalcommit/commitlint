package lint

// Failure holds Failure of linter
type Failure struct {
	inputMsg string

	failures []*RuleFailure
}

func newFailure(inputMsg string) *Failure {
	return &Failure{
		inputMsg: inputMsg,
	}
}

// AddError adds
func (res *Failure) add(r *RuleFailure) {
	res.failures = append(res.failures, r)
}

// IsOK returns true if commit message passed all the rules
func (res *Failure) IsOK() bool { return len(res.failures) == 0 }

// Input returns input commit message
func (res *Failure) Input() string { return res.inputMsg }

// Failures returns rule Failures
func (res *Failure) Failures() []*RuleFailure { return res.failures }

// RuleFailure holds Failure of a linter rule
type RuleFailure struct {
	name     string
	severity Severity
	messages []string
}

func newRuleFailure(name string, msgs []string, severity Severity) *RuleFailure {
	return &RuleFailure{
		name:     name,
		messages: msgs,
		severity: severity,
	}
}

// Name returns rule name
func (r *RuleFailure) Name() string { return r.name }

// Severity returns severity of the Rule Failure
func (r *RuleFailure) Severity() Severity { return r.severity }

// Message returns the error messages of failed rule
func (r *RuleFailure) Message() []string { return r.messages }

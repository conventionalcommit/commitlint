package commitlint

import "fmt"

// Result holds the linter error and warning messages for each rule checked
type Result struct {
	errs  []string
	warns []string
}

// NewResult returns new Result Object
func NewResult() *Result { return &Result{} }

func (res *Result) add(typ, msg string) {
	if typ == WarnType {
		res.warns = append(res.warns, msg)
	} else if typ == ErrorType {
		res.errs = append(res.errs, msg)
	} else {
		// default considered as error
		fmt.Printf("Unknown Type: %s, considering it as error", typ)
		res.errs = append(res.errs, msg)
	}
}

// Errors returns all error messages
func (res *Result) Errors() []string { return res.errs }

// Warns returns all warning messages
func (res *Result) Warns() []string { return res.warns }

// HasErrors returns true if errors found by linter
func (res *Result) HasErrors() bool { return len(res.errs) != 0 }

// HasWarns returns true if warnings found by linter
func (res *Result) HasWarns() bool { return len(res.warns) != 0 }

// IsOK returns true if commit message passed all the rules
func (res *Result) IsOK() bool { return !res.HasErrors() && !res.HasWarns() }

package cmd

import (
	"os"
	"strings"
)

// Run runs commitlint cli with os.Args
func Run() error {
	err := newCliApp().Run(os.Args)
	return handleError(err, "Failed to run commitlint CLI")
}

type multiError []error

func (m *multiError) Error() string {
	errs := make([]string, len(*m))
	for i, err := range *m {
		errs[i] = err.Error()
	}
	return strings.Join(errs, "\n")
}

func (m *multiError) Errors() []error {
	errs := make([]error, len(*m))
	for _, err := range *m {
		errs = append(errs, err)
	}
	return errs
}

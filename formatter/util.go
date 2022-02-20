package formatter

import "github.com/conventionalcommit/commitlint/lint"

// bySeverity returns all messages with given severity
func bySeverity(res *lint.Result) (errs, warns, others []*lint.Issue) {
	for _, r := range res.Issues() {
		switch r.Severity() {
		case lint.SeverityError:
			errs = append(errs, r)
		case lint.SeverityWarn:
			warns = append(warns, r)
		default:
			others = append(others, r)
		}
	}
	return errs, warns, others
}

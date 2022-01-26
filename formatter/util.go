package formatter

import "github.com/conventionalcommit/commitlint/lint"

// bySeverity returns all messages with given severity
func bySeverity(res *lint.Failure) (errs, warns, others []*lint.RuleFailure) {
	for _, r := range res.Failures() {
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

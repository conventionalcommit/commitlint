package rule

import "github.com/conventionalcommit/commitlint/lint"

// Compile-time interface checks
var (
	_ lint.Rule = (*HeadMinLenRule)(nil)
	_ lint.Rule = (*HeadMaxLenRule)(nil)
	_ lint.Rule = (*BodyMinLenRule)(nil)
	_ lint.Rule = (*BodyMaxLenRule)(nil)
	_ lint.Rule = (*BodyMaxLineLenRule)(nil)
	_ lint.Rule = (*FooterMinLenRule)(nil)
	_ lint.Rule = (*FooterMaxLenRule)(nil)
	_ lint.Rule = (*FooterMaxLineLenRule)(nil)
	_ lint.Rule = (*TypeMinLenRule)(nil)
	_ lint.Rule = (*TypeMaxLenRule)(nil)
	_ lint.Rule = (*ScopeMinLenRule)(nil)
	_ lint.Rule = (*ScopeMaxLenRule)(nil)
	_ lint.Rule = (*DescriptionMinLenRule)(nil)
	_ lint.Rule = (*DescriptionMaxLenRule)(nil)
)

// applyIntArg is a shared helper that extracts an int argument from a RuleSetting.
func applyIntArg(dst *int, ruleName string, setting lint.RuleSetting) error {
	if err := setIntArg(dst, setting.Argument); err != nil {
		return errInvalidArg(ruleName, err)
	}
	return nil
}

// --- Header ---

// HeadMinLenRule to validate min length of header
type HeadMinLenRule struct{ CheckLen int }

func (r *HeadMinLenRule) Name() string { return "header-min-length" }
func (r *HeadMinLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *HeadMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("header", r.CheckLen, msg.Header())
}

// HeadMaxLenRule to validate max length of header
type HeadMaxLenRule struct{ CheckLen int }

func (r *HeadMaxLenRule) Name() string { return "header-max-length" }
func (r *HeadMaxLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *HeadMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("header", r.CheckLen, msg.Header())
}

// --- Body ---

// BodyMinLenRule to validate min length of body
type BodyMinLenRule struct{ CheckLen int }

func (r *BodyMinLenRule) Name() string { return "body-min-length" }
func (r *BodyMinLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *BodyMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("body", r.CheckLen, msg.Body())
}

// BodyMaxLenRule to validate max length of body
type BodyMaxLenRule struct{ CheckLen int }

func (r *BodyMaxLenRule) Name() string { return "body-max-length" }
func (r *BodyMaxLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *BodyMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("body", r.CheckLen, msg.Body())
}

// BodyMaxLineLenRule to validate max line length of body
type BodyMaxLineLenRule struct{ CheckLen int }

func (r *BodyMaxLineLenRule) Name() string { return "body-max-line-length" }
func (r *BodyMaxLineLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *BodyMaxLineLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLineLength("body", r.CheckLen, msg.Body())
}

// --- Footer ---

// FooterMinLenRule to validate min length of footer
type FooterMinLenRule struct{ CheckLen int }

func (r *FooterMinLenRule) Name() string { return "footer-min-length" }
func (r *FooterMinLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *FooterMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("footer", r.CheckLen, msg.Footer())
}

// FooterMaxLenRule to validate max length of footer
type FooterMaxLenRule struct{ CheckLen int }

func (r *FooterMaxLenRule) Name() string { return "footer-max-length" }
func (r *FooterMaxLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *FooterMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("footer", r.CheckLen, msg.Footer())
}

// FooterMaxLineLenRule to validate max line length of footer
type FooterMaxLineLenRule struct{ CheckLen int }

func (r *FooterMaxLineLenRule) Name() string { return "footer-max-line-length" }
func (r *FooterMaxLineLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *FooterMaxLineLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLineLength("footer", r.CheckLen, msg.Footer())
}

// --- Type ---

// TypeMinLenRule to validate min length of type
type TypeMinLenRule struct{ CheckLen int }

func (r *TypeMinLenRule) Name() string { return "type-min-length" }
func (r *TypeMinLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *TypeMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("type", r.CheckLen, msg.Type())
}

// TypeMaxLenRule to validate max length of type
type TypeMaxLenRule struct{ CheckLen int }

func (r *TypeMaxLenRule) Name() string { return "type-max-length" }
func (r *TypeMaxLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *TypeMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("type", r.CheckLen, msg.Type())
}

// --- Scope ---

// ScopeMinLenRule to validate min length of scope
type ScopeMinLenRule struct{ CheckLen int }

func (r *ScopeMinLenRule) Name() string { return "scope-min-length" }
func (r *ScopeMinLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *ScopeMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("scope", r.CheckLen, msg.Scope())
}

// ScopeMaxLenRule to validate max length of scope
type ScopeMaxLenRule struct{ CheckLen int }

func (r *ScopeMaxLenRule) Name() string { return "scope-max-length" }
func (r *ScopeMaxLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *ScopeMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("scope", r.CheckLen, msg.Scope())
}

// --- Description ---

// DescriptionMinLenRule to validate min length of description
type DescriptionMinLenRule struct{ CheckLen int }

func (r *DescriptionMinLenRule) Name() string { return "description-min-length" }
func (r *DescriptionMinLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *DescriptionMinLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMinLen("description", r.CheckLen, msg.Description())
}

// DescriptionMaxLenRule to validate max length of description
type DescriptionMaxLenRule struct{ CheckLen int }

func (r *DescriptionMaxLenRule) Name() string { return "description-max-length" }
func (r *DescriptionMaxLenRule) Apply(s lint.RuleSetting) error {
	return applyIntArg(&r.CheckLen, r.Name(), s)
}

func (r *DescriptionMaxLenRule) Validate(msg lint.Commit) (*lint.Issue, bool) {
	return validateMaxLen("description", r.CheckLen, msg.Description())
}

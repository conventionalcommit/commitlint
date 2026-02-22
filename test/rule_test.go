package test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/internal/casing"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/rule"
)

// --- Header length rules ---

func TestHeadMinLen_Pass(t *testing.T) {
	r := &rule.HeadMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 10}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: add new feature"})
	if !ok {
		t.Error("header length >= 10 should pass")
	}
}

func TestHeadMinLen_Exact(t *testing.T) {
	r := &rule.HeadMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 10}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "0123456789"})
	if !ok {
		t.Error("header length == 10 should pass")
	}
}

func TestHeadMinLen_Fail(t *testing.T) {
	r := &rule.HeadMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 10}); err != nil {
		t.Fatal(err)
	}
	issue, ok := r.Validate(&mockCommit{header: "short"})
	if ok {
		t.Error("header length < 10 should fail")
	}
	if issue == nil || issue.Description() == "" {
		t.Error("expected non-empty issue description")
	}
}

func TestHeadMinLen_BadArg(t *testing.T) {
	r := &rule.HeadMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "string"}); err == nil {
		t.Error("expected error for non-int arg")
	}
}

func TestHeadMaxLen_Pass(t *testing.T) {
	r := &rule.HeadMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 50}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: short"})
	if !ok {
		t.Error("header <= 50 should pass")
	}
}

func TestHeadMaxLen_Exact(t *testing.T) {
	r := &rule.HeadMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 10}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "0123456789"})
	if !ok {
		t.Error("header == max should pass")
	}
}

func TestHeadMaxLen_Fail(t *testing.T) {
	r := &rule.HeadMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 10}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "this header is way too long"})
	if ok {
		t.Error("header > 10 should fail")
	}
}

func TestHeadMaxLen_NegativeDisables(t *testing.T) {
	r := &rule.HeadMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: -1}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "extremely long header that normally fails"})
	if !ok {
		t.Error("max=-1 should disable check")
	}
}

// --- Body length rules ---

func TestBodyMinLen_ZeroAllowsEmpty(t *testing.T) {
	r := &rule.BodyMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 0}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: ""})
	if !ok {
		t.Error("empty body with min=0 should pass")
	}
}

func TestBodyMinLen_Fail(t *testing.T) {
	r := &rule.BodyMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 20}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "short"})
	if ok {
		t.Error("body < 20 should fail")
	}
}

func TestBodyMaxLen_NegativeDisables(t *testing.T) {
	r := &rule.BodyMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: -1}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "very long body text here that would normally fail"})
	if !ok {
		t.Error("max=-1 should disable check")
	}
}

func TestBodyMaxLen_Fail(t *testing.T) {
	r := &rule.BodyMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 5}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "this body is too long"})
	if ok {
		t.Error("body > 5 should fail")
	}
}

func TestBodyMaxLineLen_Pass(t *testing.T) {
	r := &rule.BodyMaxLineLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 72}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "Short line 1.\nShort line 2."})
	if !ok {
		t.Error("all lines <= 72 should pass")
	}
}

func TestBodyMaxLineLen_Fail(t *testing.T) {
	r := &rule.BodyMaxLineLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 20}); err != nil {
		t.Fatal(err)
	}
	issue, ok := r.Validate(&mockCommit{body: "Short.\nThis line is definitely longer than twenty characters."})
	if ok {
		t.Error("line > 20 should fail")
	}
	if issue == nil || len(issue.Infos()) == 0 {
		t.Error("expected infos with per-line detail")
	}
}

func TestBodyMaxLineLen_EmptyBody(t *testing.T) {
	r := &rule.BodyMaxLineLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 72}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: ""})
	if !ok {
		t.Error("empty body should pass")
	}
}

// --- Footer length rules ---

func TestFooterMinLen_ZeroAllowsEmpty(t *testing.T) {
	r := &rule.FooterMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 0}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{footer: ""})
	if !ok {
		t.Error("empty footer with min=0 should pass")
	}
}

func TestFooterMaxLen_NegativeDisables(t *testing.T) {
	r := &rule.FooterMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: -1}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{footer: "Very long footer text"})
	if !ok {
		t.Error("max=-1 should disable check")
	}
}

func TestFooterMaxLineLen_Pass(t *testing.T) {
	r := &rule.FooterMaxLineLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 72}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{footer: "Fixes: #123\nReviewed-by: John"})
	if !ok {
		t.Error("footer lines <= 72 should pass")
	}
}

func TestFooterMaxLineLen_Fail(t *testing.T) {
	r := &rule.FooterMaxLineLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 10}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{footer: "Fixes: #123456789012345"})
	if ok {
		t.Error("footer line > 10 should fail")
	}
}

// --- Type rules ---

func TestTypeEnum_ValidTypes(t *testing.T) {
	r := &rule.TypeEnumRule{}
	if err := r.Apply(lint.RuleSetting{Argument: []interface{}{"feat", "fix", "docs"}}); err != nil {
		t.Fatal(err)
	}
	for _, typ := range []string{"feat", "fix", "docs"} {
		_, ok := r.Validate(&mockCommit{typ: typ})
		if !ok {
			t.Errorf("type %q should be valid", typ)
		}
	}
}

func TestTypeEnum_Invalid(t *testing.T) {
	r := &rule.TypeEnumRule{}
	if err := r.Apply(lint.RuleSetting{Argument: []interface{}{"feat", "fix"}}); err != nil {
		t.Fatal(err)
	}
	issue, ok := r.Validate(&mockCommit{typ: "unknown"})
	if ok {
		t.Error("type 'unknown' should fail")
	}
	if issue == nil {
		t.Fatal("expected non-nil issue")
	}
}

func TestTypeEnum_BadArg(t *testing.T) {
	r := &rule.TypeEnumRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "not-an-array"}); err == nil {
		t.Error("expected error for non-array arg")
	}
}

func TestTypeMinLen_Pass(t *testing.T) {
	r := &rule.TypeMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 3}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "feat"})
	if !ok {
		t.Error("type 'feat' len >= 3 should pass")
	}
}

func TestTypeMinLen_Fail(t *testing.T) {
	r := &rule.TypeMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 5}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "fix"})
	if ok {
		t.Error("type 'fix' len < 5 should fail")
	}
}

func TestTypeMaxLen_Pass(t *testing.T) {
	r := &rule.TypeMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 10}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "feat"})
	if !ok {
		t.Error("type 'feat' len <= 10 should pass")
	}
}

func TestTypeMaxLen_Fail(t *testing.T) {
	r := &rule.TypeMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 3}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "refactor"})
	if ok {
		t.Error("type 'refactor' len > 3 should fail")
	}
}

func TestTypeCharset_Pass(t *testing.T) {
	r := &rule.TypeCharsetRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "feat"})
	if !ok {
		t.Error("lowercase type should pass")
	}
}

func TestTypeCharset_Fail(t *testing.T) {
	r := &rule.TypeCharsetRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "FEAT"})
	if ok {
		t.Error("uppercase type should fail")
	}
}

func TestTypeCharset_BadArg(t *testing.T) {
	r := &rule.TypeCharsetRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 123}); err == nil {
		t.Error("expected error for non-string arg")
	}
}

// --- Scope rules ---

func TestScopeEnum_Valid(t *testing.T) {
	r := &rule.ScopeEnumRule{}
	if err := r.Apply(lint.RuleSetting{
		Argument: []interface{}{"auth", "core", "api"},
		Flags:    map[string]interface{}{"allow-empty": false},
	}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "auth"})
	if !ok {
		t.Error("scope 'auth' should pass")
	}
}

func TestScopeEnum_Invalid(t *testing.T) {
	r := &rule.ScopeEnumRule{}
	if err := r.Apply(lint.RuleSetting{
		Argument: []interface{}{"auth", "core"},
		Flags:    map[string]interface{}{"allow-empty": false},
	}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "unknown"})
	if ok {
		t.Error("scope 'unknown' should fail")
	}
}

func TestScopeEnum_EmptyAllowed(t *testing.T) {
	r := &rule.ScopeEnumRule{}
	if err := r.Apply(lint.RuleSetting{
		Argument: []interface{}{"auth"},
		Flags:    map[string]interface{}{"allow-empty": true},
	}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: ""})
	if !ok {
		t.Error("empty scope with allow-empty=true should pass")
	}
}

func TestScopeEnum_EmptyNotAllowed(t *testing.T) {
	r := &rule.ScopeEnumRule{}
	if err := r.Apply(lint.RuleSetting{
		Argument: []interface{}{"auth"},
		Flags:    map[string]interface{}{"allow-empty": false},
	}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: ""})
	if ok {
		t.Error("empty scope with allow-empty=false should fail")
	}
}

func TestScopeMinLen_Pass(t *testing.T) {
	r := &rule.ScopeMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 2}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "auth"})
	if !ok {
		t.Error("scope len >= 2 should pass")
	}
}

func TestScopeMaxLen_NegativeDisables(t *testing.T) {
	r := &rule.ScopeMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: -1}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "very-long-scope-name"})
	if !ok {
		t.Error("max=-1 should disable check")
	}
}

func TestScopeCharset_Pass(t *testing.T) {
	r := &rule.ScopeCharsetRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "abcdefghijklmnopqrstuvwxyz/"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "auth/core"})
	if !ok {
		t.Error("scope with / should pass")
	}
}

func TestScopeCharset_Fail(t *testing.T) {
	r := &rule.ScopeCharsetRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "abcdefghijklmnopqrstuvwxyz"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "Auth"})
	if ok {
		t.Error("uppercase scope should fail")
	}
}

// --- Description length rules ---

func TestDescMinLen_Pass(t *testing.T) {
	r := &rule.DescriptionMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 5}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "add new feature"})
	if !ok {
		t.Error("desc len >= 5 should pass")
	}
}

func TestDescMinLen_Fail(t *testing.T) {
	r := &rule.DescriptionMinLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 20}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "short"})
	if ok {
		t.Error("desc len < 20 should fail")
	}
}

func TestDescMaxLen_NegativeDisables(t *testing.T) {
	r := &rule.DescriptionMaxLenRule{}
	if err := r.Apply(lint.RuleSetting{Argument: -1}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "very long description text"})
	if !ok {
		t.Error("max=-1 should disable check")
	}
}

// --- Footer enum rule ---

func TestFooterEnum_Valid(t *testing.T) {
	r := &rule.FooterEnumRule{}
	if err := r.Apply(lint.RuleSetting{
		Argument: []interface{}{"Fixes", "Reviewed-by", "BREAKING CHANGE"},
	}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{
		notes: []lint.Note{&mockNote{token: "Fixes", value: "#123"}},
	})
	if !ok {
		t.Error("known footer token should pass")
	}
}

func TestFooterEnum_Invalid(t *testing.T) {
	r := &rule.FooterEnumRule{}
	if err := r.Apply(lint.RuleSetting{Argument: []interface{}{"Fixes"}}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{
		notes: []lint.Note{&mockNote{token: "Unknown-Token", value: "val"}},
	})
	if ok {
		t.Error("unknown footer token should fail")
	}
}

func TestFooterEnum_EmptyTokenList(t *testing.T) {
	r := &rule.FooterEnumRule{}
	if err := r.Apply(lint.RuleSetting{Argument: []interface{}{}}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{notes: []lint.Note{}})
	if !ok {
		t.Error("no notes with empty list should pass")
	}
}

// --- Footer type enum rule ---

func TestFooterTypeEnum_BadArg(t *testing.T) {
	r := &rule.FooterTypeEnumRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "not-an-array"}); err == nil {
		t.Error("expected error for non-array arg")
	}
}

func TestFooterTypeEnum_EmptyParams(t *testing.T) {
	r := &rule.FooterTypeEnumRule{}
	if err := r.Apply(lint.RuleSetting{Argument: []interface{}{}}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "feat"})
	if !ok {
		t.Error("empty params should pass")
	}
}

// --- Issue properties ---

func TestIssue_Properties(t *testing.T) {
	issue := lint.NewIssue("test desc", "info1", "info2")
	if issue.Description() != "test desc" {
		t.Errorf("got desc %q", issue.Description())
	}
	if len(issue.Infos()) != 2 {
		t.Errorf("expected 2 infos, got %d", len(issue.Infos()))
	}
	if issue.Infos()[0] != "info1" {
		t.Errorf("got info[0] %q", issue.Infos()[0])
	}
}

func TestIssue_NoInfos(t *testing.T) {
	issue := lint.NewIssue("desc only")
	if len(issue.Infos()) != 0 {
		t.Errorf("expected 0 infos, got %d", len(issue.Infos()))
	}
}

// ============================================================
// Case rules
// ============================================================

func TestTypeCaseRule_LowerPass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "feat"})
	if !ok {
		t.Error("lowercase type should pass lower-case rule")
	}
}

func TestTypeCaseRule_LowerFail(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "Feat"})
	if ok {
		t.Error("mixed-case type should fail lower-case rule")
	}
}

func TestTypeCaseRule_UpperPass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Upper}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "FEAT"})
	if !ok {
		t.Error("uppercase type should pass upper-case rule")
	}
}

func TestTypeCaseRule_UpperFail(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Upper}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "feat"})
	if ok {
		t.Error("lowercase type should fail upper-case rule")
	}
}

func TestTypeCaseRule_CamelPass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Camel}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "myType"})
	if !ok {
		t.Error("camelCase type should pass")
	}
}

func TestTypeCaseRule_PascalPass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Pascal}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "MyType"})
	if !ok {
		t.Error("PascalCase type should pass")
	}
}

func TestTypeCaseRule_KebabPass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Kebab}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "my-type"})
	if !ok {
		t.Error("kebab-case type should pass")
	}
}

func TestTypeCaseRule_SnakePass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Snake}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "my_type"})
	if !ok {
		t.Error("snake_case type should pass")
	}
}

func TestTypeCaseRule_SentencePass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Sentence}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "Feat"})
	if !ok {
		t.Error("Sentence case type should pass")
	}
}

func TestTypeCaseRule_StartPass(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Start}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "My Type"})
	if !ok {
		t.Error("Start Case type should pass")
	}
}

func TestTypeCaseRule_BadCaseArg(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "unknown-case"}); err == nil {
		t.Error("unknown case should return error")
	}
}

func TestTypeCaseRule_BadArgType(t *testing.T) {
	r := &rule.TypeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 42}); err == nil {
		t.Error("non-string arg should return error")
	}
}

func TestScopeCaseRule_LowerPass(t *testing.T) {
	r := &rule.ScopeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "auth"})
	if !ok {
		t.Error("lowercase scope should pass")
	}
}

func TestScopeCaseRule_EmptyScopeAlwaysPasses(t *testing.T) {
	r := &rule.ScopeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: ""})
	if !ok {
		t.Error("empty scope should always pass scope-case")
	}
}

func TestScopeCaseRule_Fail(t *testing.T) {
	r := &rule.ScopeCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "Auth"})
	if ok {
		t.Error("uppercase scope should fail lower-case rule")
	}
}

func TestDescriptionCaseRule_LowerPass(t *testing.T) {
	r := &rule.DescriptionCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "add new feature"})
	if !ok {
		t.Error("lowercase description should pass")
	}
}

func TestDescriptionCaseRule_Fail(t *testing.T) {
	r := &rule.DescriptionCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "Add new feature"})
	if ok {
		t.Error("capitalized description should fail lower-case rule")
	}
}

func TestBodyCaseRule_LowerPass(t *testing.T) {
	r := &rule.BodyCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "this is the body"})
	if !ok {
		t.Error("lowercase body should pass")
	}
}

func TestBodyCaseRule_EmptyBodyAlwaysPasses(t *testing.T) {
	r := &rule.BodyCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: ""})
	if !ok {
		t.Error("empty body should pass body-case")
	}
}

func TestBodyCaseRule_Fail(t *testing.T) {
	r := &rule.BodyCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "This is the body"})
	if ok {
		t.Error("capitalized body should fail lower-case rule")
	}
}

func TestBodyCaseRule_MultiLine_AllFail(t *testing.T) {
	r := &rule.BodyCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	// Entire body fails lower-case (uppercase letters present)
	_, ok := r.Validate(&mockCommit{body: "First line capitalized\nSecond line capitalized"})
	if ok {
		t.Error("body with uppercase letters should fail lower-case rule")
	}
}

func TestBodyCaseRule_MultiLine_SomeFail(t *testing.T) {
	r := &rule.BodyCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	// Whole body contains uppercase, so the entire body fails
	_, ok := r.Validate(&mockCommit{body: "first line is good\nSecond line is bad"})
	if ok {
		t.Error("body containing an uppercase letter should fail lower-case rule")
	}
}

func TestBodyCaseRule_MultiLine_BlankLineSkipped(t *testing.T) {
	r := &rule.BodyCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	// All non-empty content is lowercase; blank line is fine
	_, ok := r.Validate(&mockCommit{body: "first line\n\nsecond line"})
	if !ok {
		t.Error("all-lowercase body with blank separator should pass lower-case rule")
	}
}

func TestHeaderCaseRule_LowerPass(t *testing.T) {
	r := &rule.HeaderCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: add feature"})
	if !ok {
		t.Error("lowercase header should pass")
	}
}

func TestHeaderCaseRule_Fail(t *testing.T) {
	r := &rule.HeaderCaseRule{}
	if err := r.Apply(lint.RuleSetting{Argument: casing.Lower}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "Feat: Add feature"})
	if ok {
		t.Error("capitalized header should fail lower-case rule")
	}
}

// ============================================================
// Empty rules
// ============================================================

func TestTypeEmptyRule_NonEmptyPasses(t *testing.T) {
	r := &rule.TypeEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{typ: "feat"})
	if !ok {
		t.Error("non-empty type should pass type-empty rule")
	}
}

func TestTypeEmptyRule_EmptyFails(t *testing.T) {
	r := &rule.TypeEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	issue, ok := r.Validate(&mockCommit{typ: ""})
	if ok {
		t.Error("empty type should fail type-empty rule")
	}
	if issue == nil {
		t.Error("expected non-nil issue")
	}
}

func TestScopeEmptyRule_NonEmptyPasses(t *testing.T) {
	r := &rule.ScopeEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: "auth"})
	if !ok {
		t.Error("non-empty scope should pass scope-empty rule")
	}
}

func TestScopeEmptyRule_EmptyFails(t *testing.T) {
	r := &rule.ScopeEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{scope: ""})
	if ok {
		t.Error("empty scope should fail scope-empty rule")
	}
}

func TestBodyEmptyRule_NonEmptyPasses(t *testing.T) {
	r := &rule.BodyEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "some body"})
	if !ok {
		t.Error("non-empty body should pass body-empty rule")
	}
}

func TestBodyEmptyRule_EmptyFails(t *testing.T) {
	r := &rule.BodyEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: ""})
	if ok {
		t.Error("empty body should fail body-empty rule")
	}
}

func TestFooterEmptyRule_NonEmptyPasses(t *testing.T) {
	r := &rule.FooterEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{footer: "Fixes: #123"})
	if !ok {
		t.Error("non-empty footer should pass footer-empty rule")
	}
}

func TestFooterEmptyRule_EmptyFails(t *testing.T) {
	r := &rule.FooterEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{footer: ""})
	if ok {
		t.Error("empty footer should fail footer-empty rule")
	}
}

func TestDescriptionEmptyRule_NonEmptyPasses(t *testing.T) {
	r := &rule.DescriptionEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "add new feature"})
	if !ok {
		t.Error("non-empty description should pass description-empty rule")
	}
}

func TestDescriptionEmptyRule_EmptyFails(t *testing.T) {
	r := &rule.DescriptionEmptyRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: ""})
	if ok {
		t.Error("empty description should fail description-empty rule")
	}
}

// ============================================================
// Full-stop rules
// ============================================================

func TestHeaderFullStop_NoStop_Pass(t *testing.T) {
	r := &rule.HeaderFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: add new feature"})
	if !ok {
		t.Error("header not ending with '.' should pass")
	}
}

func TestHeaderFullStop_WithStop_Fail(t *testing.T) {
	r := &rule.HeaderFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: add new feature."})
	if ok {
		t.Error("header ending with '.' should fail")
	}
}

func TestHeaderFullStop_CustomChar(t *testing.T) {
	r := &rule.HeaderFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "!"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: urgent!"})
	if ok {
		t.Error("header ending with '!' should fail")
	}
	_, ok2 := r.Validate(&mockCommit{header: "feat: normal"})
	if !ok2 {
		t.Error("header not ending with '!' should pass")
	}
}

func TestHeaderFullStop_BadArg(t *testing.T) {
	r := &rule.HeaderFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 99}); err == nil {
		t.Error("non-string arg should return error")
	}
}

func TestBodyFullStop_NoStop_Pass(t *testing.T) {
	r := &rule.BodyFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "This is the body"})
	if !ok {
		t.Error("body not ending with '.' should pass")
	}
}

func TestBodyFullStop_WithStop_Fail(t *testing.T) {
	r := &rule.BodyFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: "This is the body."})
	if ok {
		t.Error("body ending with '.' should fail")
	}
}

func TestBodyFullStop_EmptyBody_Pass(t *testing.T) {
	r := &rule.BodyFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: ""})
	if !ok {
		t.Error("empty body should pass body-full-stop")
	}
}

func TestDescriptionFullStop_NoStop_Pass(t *testing.T) {
	r := &rule.DescriptionFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "add new feature"})
	if !ok {
		t.Error("description not ending with '.' should pass")
	}
}

func TestDescriptionFullStop_WithStop_Fail(t *testing.T) {
	r := &rule.DescriptionFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: "add new feature."})
	if ok {
		t.Error("description ending with '.' should fail")
	}
}

func TestDescriptionFullStop_EmptyDescription_Pass(t *testing.T) {
	r := &rule.DescriptionFullStopRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "."}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{description: ""})
	if !ok {
		t.Error("empty description should pass description-full-stop")
	}
}

// ============================================================
// Leading-blank rules
// ============================================================

func TestBodyLeadingBlank_WithBlank_Pass(t *testing.T) {
	r := &rule.BodyLeadingBlankRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	// message: header + blank line + body
	msg := &mockCommit{
		message: "feat: add feature\n\nThis is the body",
		body:    "This is the body",
	}
	_, ok := r.Validate(msg)
	if !ok {
		t.Error("body with leading blank line should pass")
	}
}

func TestBodyLeadingBlank_WithoutBlank_Fail(t *testing.T) {
	r := &rule.BodyLeadingBlankRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		message: "feat: add feature\nThis is the body",
		body:    "This is the body",
	}
	_, ok := r.Validate(msg)
	if ok {
		t.Error("body without leading blank line should fail")
	}
}

func TestBodyLeadingBlank_EmptyBody_Pass(t *testing.T) {
	r := &rule.BodyLeadingBlankRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{body: ""})
	if !ok {
		t.Error("empty body should pass body-leading-blank")
	}
}

func TestFooterLeadingBlank_WithBlank_Pass(t *testing.T) {
	r := &rule.FooterLeadingBlankRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		message: "feat: add feature\n\nbody text\n\nFixes: #123",
		footer:  "Fixes: #123",
	}
	_, ok := r.Validate(msg)
	if !ok {
		t.Error("footer with leading blank should pass")
	}
}

func TestFooterLeadingBlank_WithoutBlank_Fail(t *testing.T) {
	r := &rule.FooterLeadingBlankRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		message: "feat: add feature\nbody text\nFixes: #123",
		footer:  "Fixes: #123",
	}
	_, ok := r.Validate(msg)
	if ok {
		t.Error("footer without leading blank should fail")
	}
}

func TestFooterLeadingBlank_EmptyFooter_Pass(t *testing.T) {
	r := &rule.FooterLeadingBlankRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{footer: ""})
	if !ok {
		t.Error("empty footer should pass footer-leading-blank")
	}
}

// ============================================================
// Header-trim rule
// ============================================================

func TestHeaderTrim_Clean_Pass(t *testing.T) {
	r := &rule.HeaderTrimRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: add feature"})
	if !ok {
		t.Error("clean header should pass header-trim")
	}
}

func TestHeaderTrim_LeadingSpace_Fail(t *testing.T) {
	r := &rule.HeaderTrimRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: " feat: add feature"})
	if ok {
		t.Error("header with leading space should fail header-trim")
	}
}

func TestHeaderTrim_TrailingSpace_Fail(t *testing.T) {
	r := &rule.HeaderTrimRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "feat: add feature "})
	if ok {
		t.Error("header with trailing space should fail header-trim")
	}
}

func TestHeaderTrim_BothSpaces_Fail(t *testing.T) {
	r := &rule.HeaderTrimRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{header: "  feat: add feature  "})
	if ok {
		t.Error("header with both-sides whitespace should fail header-trim")
	}
}

// ============================================================
// Signed-off-by and trailer-exists rules
// ============================================================

func TestSignedOffBy_Present_Pass(t *testing.T) {
	r := &rule.SignedOffByRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "Signed-off-by:"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{notes: []lint.Note{&mockNote{token: "Signed-off-by", value: "Jane Doe <jane@example.com>"}}})
	if !ok {
		t.Error("message with Signed-off-by should pass")
	}
}

func TestSignedOffBy_NoColon_Present_Pass(t *testing.T) {
	r := &rule.SignedOffByRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "Signed-off-by"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{notes: []lint.Note{&mockNote{token: "Signed-off-by", value: "Jane Doe <jane@example.com>"}}})
	if !ok {
		t.Error("message with Signed-off-by should pass")
	}
}

func TestSignedOffBy_Missing_Fail(t *testing.T) {
	r := &rule.SignedOffByRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "Signed-off-by:"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{})
	if ok {
		t.Error("commit without Signed-off-by note should fail")
	}
}

func TestSignedOffBy_NoColon_Missing_Fail(t *testing.T) {
	r := &rule.SignedOffByRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "Signed-off-by"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{})
	if ok {
		t.Error("commit without Signed-off-by note should fail")
	}
}

func TestSignedOffBy_BadArg(t *testing.T) {
	r := &rule.SignedOffByRule{}
	if err := r.Apply(lint.RuleSetting{Argument: 99}); err == nil {
		t.Error("non-string arg should return error")
	}
}

func TestTrailerExists_Present_Pass(t *testing.T) {
	r := &rule.TrailerExistsRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "Co-authored-by:"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{notes: []lint.Note{&mockNote{token: "Co-authored-by", value: "Bob <bob@example.com>"}}})
	if !ok {
		t.Error("commit with Co-authored-by note should pass")
	}
}

func TestTrailerExists_Missing_Fail(t *testing.T) {
	r := &rule.TrailerExistsRule{}
	if err := r.Apply(lint.RuleSetting{Argument: "Co-authored-by:"}); err != nil {
		t.Fatal(err)
	}
	_, ok := r.Validate(&mockCommit{})
	if ok {
		t.Error("commit without Co-authored-by note should fail")
	}
}

// ============================================================
// Breaking-change-exclamation-mark rule
// ============================================================

func TestBreakingChangeExclamation_BothPresent_Pass(t *testing.T) {
	r := &rule.BreakingChangeExclamationMarkRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		breaking: true,
		notes:    []lint.Note{&mockNote{token: "BREAKING CHANGE", value: "removed old endpoint"}},
	}
	_, ok := r.Validate(msg)
	if !ok {
		t.Error("both '!' in header AND BREAKING CHANGE in footer should pass")
	}
}

func TestBreakingChangeExclamation_NeitherPresent_Pass(t *testing.T) {
	r := &rule.BreakingChangeExclamationMarkRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		breaking: false,
	}
	_, ok := r.Validate(msg)
	if !ok {
		t.Error("neither '!' nor BREAKING CHANGE should pass (XNOR)")
	}
}

func TestBreakingChangeExclamation_OnlyExclamation_Fail(t *testing.T) {
	r := &rule.BreakingChangeExclamationMarkRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		breaking: true,
	}
	_, ok := r.Validate(msg)
	if ok {
		t.Error("'!' in header without BREAKING CHANGE in footer should fail")
	}
}

func TestBreakingChangeExclamation_OnlyFooter_Fail(t *testing.T) {
	r := &rule.BreakingChangeExclamationMarkRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		breaking: false,
		notes:    []lint.Note{&mockNote{token: "BREAKING CHANGE", value: "some change"}},
	}
	_, ok := r.Validate(msg)
	if ok {
		t.Error("BREAKING CHANGE in footer without '!' in header should fail")
	}
}

func TestBreakingChangeExclamation_BreakingDashChange_Pass(t *testing.T) {
	r := &rule.BreakingChangeExclamationMarkRule{}
	if err := r.Apply(lint.RuleSetting{}); err != nil {
		t.Fatal(err)
	}
	msg := &mockCommit{
		breaking: true,
		notes:    []lint.Note{&mockNote{token: "BREAKING-CHANGE", value: "some change"}},
	}
	_, ok := r.Validate(msg)
	if !ok {
		t.Error("BREAKING-CHANGE (with dash) in footer with '!' should pass")
	}
}

package test

import (
	"testing"

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

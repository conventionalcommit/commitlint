package casing_test

import (
	"testing"

	"github.com/conventionalcommit/commitlint/internal/casing"
)

func TestIsCamelCase(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"", true},            // empty always passes
		{"feat", true},        // all-lowercase single word
		{"myFeature", true},   // classic camelCase
		{"parseHTML", true},   // acronym tail is fine (still no separator)
		{"myFeat123", true},   // digits allowed
		{"MyFeature", false},  // starts uppercase → PascalCase, not camelCase
		{"my-feature", false}, // hyphen not allowed
		{"my_feature", false}, // underscore not allowed
		{"my feature", false}, // space not allowed
	} {
		got := casing.IsCamelCase(tc.in)
		if got != tc.want {
			t.Errorf("IsCamelCase(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsKebabCase(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"", true},
		{"feat", true},
		{"my-feature", true},
		{"kebab-case", true},
		{"v2-api", true},      // digit allowed
		{"my--feat", true},    // consecutive hyphens: spec doesn't forbid them
		{"MyFeature", false},  // uppercase
		{"my_feature", false}, // underscore
		{"my feature", false}, // space
		{"MY-FEAT", false},    // uppercase letters
	} {
		got := casing.IsKebabCase(tc.in)
		if got != tc.want {
			t.Errorf("IsKebabCase(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsPascalCase(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"", true},
		{"Feat", true},
		{"MyFeature", true},
		{"ParseHTML", true},
		{"MyFeat123", true},   // digits ok
		{"myFeature", false},  // starts lowercase
		{"My-Feature", false}, // hyphen not allowed
		{"My Feature", false}, // space not allowed
	} {
		got := casing.IsPascalCase(tc.in)
		if got != tc.want {
			t.Errorf("IsPascalCase(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsSentenceCase(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"", true},
		{"Feat", true},
		{"My feature", true},
		{"Add endpoint", true},
		{"Fix #123", true},    // punctuation and digit allowed
		{"feat", false},       // starts lowercase
		{"My Feature", false}, // second word capitalised
		{"MY FEATURE", false}, // fully uppercased
	} {
		got := casing.IsSentenceCase(tc.in)
		if got != tc.want {
			t.Errorf("IsSentenceCase(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsSnakeCase(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"", true},
		{"feat", true},
		{"my_feature", true},
		{"snake_case", true},
		{"v2_api", true},      // digit allowed
		{"MyFeature", false},  // uppercase
		{"my-feature", false}, // hyphen
		{"my feature", false}, // space
	} {
		got := casing.IsSnakeCase(tc.in)
		if got != tc.want {
			t.Errorf("IsSnakeCase(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsStartCase(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"", true},
		{"Feat", true},
		{"My Feature", true},
		{"Add New Endpoint", true},
		{"my feature", false}, // first word lowercase
		{"My feature", false}, // second word lowercase
		{"MY FEATURE", true},  // all-uppercase still starts uppercase per word
	} {
		got := casing.IsStartCase(tc.in)
		if got != tc.want {
			t.Errorf("IsStartCase(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestCheck_AllFormats(t *testing.T) {
	for _, tc := range []struct {
		format string
		pass   []string
		fail   []string
	}{
		{
			casing.Lower,
			[]string{"", "feat", "my feat"},
			[]string{"Feat", "FEAT"},
		},
		{
			casing.Upper,
			[]string{"", "FEAT", "MY FEAT"},
			[]string{"feat", "Feat"},
		},
		{
			casing.Camel,
			[]string{"", "feat", "myFeature"},
			[]string{"MyFeature", "my-feature"},
		},
		{
			casing.Kebab,
			[]string{"", "feat", "my-feature"},
			[]string{"MyFeature", "my_feature"},
		},
		{
			casing.Pascal,
			[]string{"", "Feat", "MyFeature"},
			[]string{"feat", "my_feature"},
		},
		{
			casing.Sentence,
			[]string{"", "Feat", "My feat"},
			[]string{"feat", "My Feat"},
		},
		{
			casing.Snake,
			[]string{"", "feat", "my_feature"},
			[]string{"MyFeature", "my-feature"},
		},
		{
			casing.Start,
			[]string{"", "Feat", "My Feature"},
			[]string{"feat", "my feature"},
		},
	} {
		for _, s := range tc.pass {
			if !casing.Check(s, tc.format) {
				t.Errorf("Check(%q, %q) = false, want true", s, tc.format)
			}
		}
		for _, s := range tc.fail {
			if casing.Check(s, tc.format) {
				t.Errorf("Check(%q, %q) = true, want false", s, tc.format)
			}
		}
	}
}

func TestCheck_UnknownFormat(t *testing.T) {
	if casing.Check("anything", "not-a-case") {
		t.Error("Check with unknown format should return false")
	}
}

func TestAll_ContainsAllConstants(t *testing.T) {
	want := map[string]bool{
		casing.Lower:    false,
		casing.Upper:    false,
		casing.Camel:    false,
		casing.Kebab:    false,
		casing.Pascal:   false,
		casing.Sentence: false,
		casing.Snake:    false,
		casing.Start:    false,
	}
	for _, c := range casing.All {
		if _, ok := want[c]; !ok {
			t.Errorf("casing.All contains unexpected value %q", c)
		}
		want[c] = true
	}
	for k, seen := range want {
		if !seen {
			t.Errorf("casing.All is missing constant %q", k)
		}
	}
}

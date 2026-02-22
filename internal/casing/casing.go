// Package casing provides case-format constants and validators used by commit
// lint rules. Each format has a corresponding exported predicate and a shared
// [Check] dispatcher.
package casing

import (
	"strings"
	"unicode"
)

// Case-format constants. These are the only values accepted by rules that take
// a case argument (e.g. type-case, scope-case, …).
const (
	Lower    = "lower-case"    // all characters lower-cased, e.g. "feat"
	Upper    = "upper-case"    // all characters upper-cased, e.g. "FEAT"
	Camel    = "camel-case"    // starts lowercase, no separators, e.g. "myFeat"
	Kebab    = "kebab-case"    // lowercase words joined by hyphens, e.g. "my-feat"
	Pascal   = "pascal-case"   // starts uppercase, no separators, e.g. "MyFeat"
	Sentence = "sentence-case" // first letter uppercase, rest lowercase, e.g. "My feat"
	Snake    = "snake-case"    // lowercase words joined by underscores, e.g. "my_feat"
	Start    = "start-case"    // every word starts uppercase, e.g. "My Feat"
)

// All is the ordered list of all valid case formats.
var All = []string{Lower, Upper, Camel, Kebab, Pascal, Sentence, Snake, Start}

// Check returns true when s conforms to the given caseFormat constant.
// An empty string always returns true (treated as "not present").
// Returns false for any unrecognised caseFormat value.
func Check(s, caseFormat string) bool {
	if s == "" {
		return true
	}
	switch caseFormat {
	case Lower:
		return s == strings.ToLower(s)
	case Upper:
		return s == strings.ToUpper(s)
	case Camel:
		return IsCamelCase(s)
	case Kebab:
		return IsKebabCase(s)
	case Pascal:
		return IsPascalCase(s)
	case Sentence:
		return IsSentenceCase(s)
	case Snake:
		return IsSnakeCase(s)
	case Start:
		return IsStartCase(s)
	default:
		return false
	}
}

// IsCamelCase reports whether s is in camelCase format.
// camelCase: starts with a lowercase letter and contains only letters and
// digits (no separators). Examples: "feat", "myFeature", "parseHTML".
func IsCamelCase(s string) bool {
	if s == "" {
		return true
	}
	runes := []rune(s)
	if unicode.IsUpper(runes[0]) {
		return false
	}
	for _, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsKebabCase reports whether s is in kebab-case format.
// kebab-case: only lowercase letters, digits, and hyphens.
// Examples: "my-feature", "kebab-case", "v2-api".
func IsKebabCase(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '-' {
			return false
		}
	}
	return true
}

// IsPascalCase reports whether s is in PascalCase format.
// PascalCase: starts with an uppercase letter and contains only letters and
// digits (no separators). Examples: "MyFeature", "ParseHTML", "Feat".
func IsPascalCase(s string) bool {
	if s == "" {
		return true
	}
	runes := []rune(s)
	if !unicode.IsUpper(runes[0]) {
		return false
	}
	for _, r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsSentenceCase reports whether s is in Sentence case format.
// Sentence case: first rune is uppercase; every subsequent letter is lowercase.
// Digits and punctuation are allowed anywhere.
// Examples: "My feature", "Add endpoint", "Fix #123".
func IsSentenceCase(s string) bool {
	if s == "" {
		return true
	}
	runes := []rune(s)
	if !unicode.IsUpper(runes[0]) {
		return false
	}
	for _, r := range runes[1:] {
		if unicode.IsLetter(r) && unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// IsSnakeCase reports whether s is in snake_case format.
// snake_case: only lowercase letters, digits, and underscores.
// Examples: "my_feature", "snake_case", "v2_api".
func IsSnakeCase(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

// IsStartCase reports whether s is in Start Case format.
// Start Case: every whitespace-separated word starts with an uppercase letter.
// Examples: "My Feature", "Add New Endpoint".
func IsStartCase(s string) bool {
	for _, w := range strings.Fields(s) {
		runes := []rune(w)
		if !unicode.IsUpper(runes[0]) {
			return false
		}
	}
	return true
}

// Package commitlint provides a simple linter for conventional commits
package commitlint

import (
	"fmt"
	"strings"

	"github.com/conventionalcommit/parser"
)

// Linter is a simple linter for conventional commits
type Linter struct {
	conf *Config
}

// NewLinter returns a new Linter object with given conf
func NewLinter(conf *Config) *Linter {
	return &Linter{conf: conf}
}

// Lint checks the given commitMsg string against rules
func (l *Linter) Lint(commitMsg string) (lintReport string, hasError bool, err error) {
	msg, err := parser.Parse(commitMsg)
	if err != nil {
		if parser.IsHeaderErr(err) {
			res := NewResult()
			// TODO: show more information
			res.add(ErrorType, "commit header is not valid")
			return l.formReport(msg, res), true, nil
		}
		return "", false, err
	}

	lintReport, hasError = l.LintCommit(msg)
	lintReport = strings.Trim(lintReport, "\n")
	return lintReport, hasError, nil
}

// LintCommit checks the given parser.Commit for rules
func (l *Linter) LintCommit(msg *parser.Commit) (lintReport string, hasError bool) {
	res := NewResult()

	l.checkHeader(msg, res)
	l.checkBody(msg, res)
	l.checkFooter(msg, res)

	report := l.formReport(msg, res)
	return report, res.HasErrors()
}

func (l *Linter) formReport(msg *parser.Commit, res *Result) string {
	var report string
	if res.IsOK() {
		report = successMessage(msg, res)
	} else {
		report = errorMessage(msg, res)
	}
	return strings.Trim(report, "\n")
}

func (l *Linter) checkHeader(msg *parser.Commit, res *Result) {
	headConf := l.conf.Header
	head := msg.Header

	// Header Min Length Check
	if minConf := headConf.MinLength; minConf.Enabled {
		actualLen := len(head.FullHeader)
		if actualLen < minConf.Value {
			errMsg := fmt.Sprintf("header: length is %d, should have atleast %d chars", actualLen, minConf.Value)
			res.add(minConf.Type, errMsg)
		}
	}

	// Header Max Length Check
	if maxConf := headConf.MaxLength; maxConf.Enabled {
		actualLen := len(head.FullHeader)
		if actualLen > maxConf.Value {
			errMsg := fmt.Sprintf("header: length is %d, should have less than %d chars", actualLen, maxConf.Value)
			res.add(maxConf.Type, errMsg)
		}
	}

	// Type Check
	if typesConf := headConf.Types; typesConf.Enabled {
		allowedTypes := typesConf.Value
		isFound := search(allowedTypes, head.Type)
		if !isFound {
			errMsg := fmt.Sprintf("type: '%s' is not allowed, you can use one of %v", head.Type, allowedTypes)
			res.add(typesConf.Type, errMsg)
		}
	}

	// Scope Check
	if scopeConf := headConf.Scopes; scopeConf.Enabled {
		allowedScopes := scopeConf.Value
		isFound := search(allowedScopes, head.Scope)
		if !isFound {
			errMsg := fmt.Sprintf("scope: '%s' is not allowed, you can use one of %v", head.Scope, allowedScopes)
			res.add(scopeConf.Type, errMsg)
		}
	}
}

func (l *Linter) checkBody(msg *parser.Commit, res *Result) {
	bodyConf := l.conf.Body

	if !bodyConf.CanBeEmpty {
		if msg.Body == "" {
			res.add(ErrorType, "body: cannot be empty")
			return
		}
	}

	// bodyer Max Line Length Check
	if maxConf := bodyConf.MaxLineLength; maxConf.Enabled {
		bodyLines := strings.Split(msg.Body, "\n")
		for index, bodyLine := range bodyLines {
			actualLen := len(bodyLine)
			if actualLen > maxConf.Value {
				errMsg := fmt.Sprintf("body: in line %d, length is %d, should have less than %d chars", index+1, actualLen, maxConf.Value)
				res.add(maxConf.Type, errMsg)
			}
		}
	}
}

func (l *Linter) checkFooter(msg *parser.Commit, res *Result) {
	footConf := l.conf.Footer
	foot := msg.Footer

	if !footConf.CanBeEmpty {
		if foot.FullFooter == "" {
			res.add(ErrorType, "footer: cannot be empty")
			return
		}
	}

	// Footer Max Line Length Check
	if maxConf := footConf.MaxLineLength; maxConf.Enabled {
		footLines := strings.Split(foot.FullFooter, "\n")
		for index, footLine := range footLines {
			actualLen := len(footLine)
			if actualLen > maxConf.Value {
				errMsg := fmt.Sprintf("footer: in line %d, length is %d, should have less than %d chars", index+1, actualLen, maxConf.Value)
				res.add(maxConf.Type, errMsg)
			}
		}
	}
}

func search(arr []string, toFind string) bool {
	for _, typ := range arr {
		if typ == toFind {
			return true
		}
	}
	return false
}

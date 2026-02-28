// Package formatter contains changelog output formatters
package formatter

import (
	"fmt"
	"strings"

	"github.com/conventionalcommit/commitlint/changelog"
)

// MarkdownFormatter formats changelog as Markdown
type MarkdownFormatter struct{}

// Name returns the name of the formatter
func (f *MarkdownFormatter) Name() string { return "markdown" }

// Format formats the full changelog as Markdown
func (f *MarkdownFormatter) Format(cl *changelog.Changelog) (string, error) {
	var sb strings.Builder

	if cl.Header != "" {
		sb.WriteString(cl.Header)
		sb.WriteString("\n")
	}

	for i, v := range cl.Versions {
		if i > 0 || cl.Header != "" {
			sb.WriteString("\n")
		}

		out, err := f.FormatVersion(&v, nil)
		if err != nil {
			return "", err
		}
		sb.WriteString(out)
	}

	return sb.String(), nil
}

// FormatVersion formats a single version changelog as Markdown
func (f *MarkdownFormatter) FormatVersion(v *changelog.VersionChangelog, _ *changelog.Config) (string, error) {
	var sb strings.Builder

	// version header
	f.writeVersionHeader(&sb, v)

	// breaking changes
	if len(v.Breaking) > 0 {
		sb.WriteString("\n### ⚠ BREAKING CHANGES\n\n")
		for _, c := range v.Breaking {
			f.writeCommitLine(&sb, c)
		}
	}

	// type groups
	for _, g := range v.Groups {
		sb.WriteString("\n### ")
		sb.WriteString(g.Header)
		sb.WriteString("\n\n")

		for _, c := range g.Commits {
			f.writeCommitLine(&sb, c)
		}
	}

	// other (non-conventional) commits
	if len(v.Other) > 0 {
		sb.WriteString("\n### Other Changes\n\n")
		for _, c := range v.Other {
			f.writeOtherCommitLine(&sb, c)
		}
	}

	return sb.String(), nil
}

// writeVersionHeader writes the version header line
func (f *MarkdownFormatter) writeVersionHeader(sb *strings.Builder, v *changelog.VersionChangelog) {
	sb.WriteString("## ")

	if v.CompareURL != "" {
		fmt.Fprintf(sb, "[%s](%s)", v.Version, v.CompareURL)
	} else {
		sb.WriteString(v.Version)
	}

	if v.Date != "" {
		fmt.Fprintf(sb, " (%s)", v.Date)
	}

	sb.WriteString("\n")
}

// writeCommitLine writes a single commit line
func (f *MarkdownFormatter) writeCommitLine(sb *strings.Builder, c changelog.CommitInfo) {
	sb.WriteString("* ")

	if c.Commit != nil && c.Commit.Scope() != "" {
		fmt.Fprintf(sb, "**%s:** ", c.Commit.Scope())
	}

	if c.Commit != nil {
		sb.WriteString(c.Commit.Description())
	}

	f.writeHashLink(sb, c)
	sb.WriteString("\n")
}

// writeOtherCommitLine writes a non-conventional commit line
func (f *MarkdownFormatter) writeOtherCommitLine(sb *strings.Builder, c changelog.CommitInfo) {
	sb.WriteString("* ")

	if c.Commit != nil {
		sb.WriteString(c.Commit.Description())
	} else {
		sb.WriteString(c.ShortHash)
	}

	f.writeHashLink(sb, c)
	sb.WriteString("\n")
}

// writeHashLink writes the commit hash with optional link
func (f *MarkdownFormatter) writeHashLink(sb *strings.Builder, c changelog.CommitInfo) {
	if c.CommitURL != "" {
		fmt.Fprintf(sb, " ([%s](%s))", c.ShortHash, c.CommitURL)
	} else if c.ShortHash != "" {
		fmt.Fprintf(sb, " (%s)", c.ShortHash)
	}
}

package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/conventionalcommit/commitlint/changelog"
)

// JSONFormatter formats changelog as JSON
type JSONFormatter struct{}

// Name returns the name of the formatter
func (f *JSONFormatter) Name() string { return "json" }

// Format formats the full changelog as JSON
func (f *JSONFormatter) Format(cl *changelog.Changelog) (string, error) {
	output := map[string]interface{}{
		"header":   cl.Header,
		"versions": f.formatVersions(cl.Versions),
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json formatting failed: %w", err)
	}
	return string(data), nil
}

// FormatVersion formats a single version changelog as JSON
func (f *JSONFormatter) FormatVersion(v *changelog.VersionChangelog, _ *changelog.Config) (string, error) {
	output := f.formatSingleVersion(v)

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json formatting failed: %w", err)
	}
	return string(data), nil
}

func (f *JSONFormatter) formatVersions(versions []changelog.VersionChangelog) []interface{} {
	result := make([]interface{}, 0, len(versions))
	for i := range versions {
		result = append(result, f.formatSingleVersion(&versions[i]))
	}
	return result
}

func (f *JSONFormatter) formatSingleVersion(v *changelog.VersionChangelog) map[string]interface{} {
	output := map[string]interface{}{
		"version": v.Version,
		"date":    v.Date,
	}

	if v.CompareURL != "" {
		output["compareUrl"] = v.CompareURL
	}

	if len(v.Breaking) > 0 {
		output["breakingChanges"] = f.formatCommits(v.Breaking)
	}

	groups := make([]interface{}, 0, len(v.Groups))
	for _, g := range v.Groups {
		groups = append(groups, map[string]interface{}{
			"type":    g.Type,
			"header":  g.Header,
			"commits": f.formatCommits(g.Commits),
		})
	}
	output["groups"] = groups

	if len(v.Other) > 0 {
		output["other"] = f.formatCommits(v.Other)
	}

	return output
}

func (f *JSONFormatter) formatCommits(commits []changelog.CommitInfo) []interface{} {
	result := make([]interface{}, 0, len(commits))
	for _, c := range commits {
		entry := map[string]interface{}{
			"hash":      c.Hash,
			"shortHash": c.ShortHash,
			"author":    c.Author,
			"date":      c.Date.Format("2006-01-02"),
		}

		if c.CommitURL != "" {
			entry["commitUrl"] = c.CommitURL
		}

		if c.Commit != nil {
			entry["type"] = c.Commit.Type()
			entry["scope"] = c.Commit.Scope()
			entry["description"] = c.Commit.Description()
			entry["isBreakingChange"] = c.Commit.IsBreakingChange()
		}

		if len(c.References) > 0 {
			entry["references"] = c.References
		}

		result = append(result, entry)
	}
	return result
}

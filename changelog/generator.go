package changelog

import (
	"fmt"
	"strings"
	"time"

	"github.com/conventionalcommit/commitlint/commit"
	"github.com/conventionalcommit/commitlint/internal/git"
)

// Generator orchestrates changelog generation
type Generator struct {
	git  *git.Client
	conf *Config

	parser     commit.Parser
	commitURL  string
	compareURL string
}

// New creates a new Generator for the given repository directory
func New(repoDir string, conf *Config) (*Generator, error) {
	git, err := git.NewClient(repoDir)
	if err != nil {
		return nil, err
	}

	g := &Generator{
		git:    git,
		conf:   conf,
		parser: commit.NewParser(),
	}

	// infer repository URLs if not configured
	if err := g.inferURLs(); err != nil {
		return nil, err
	}

	return g, nil
}

// inferURLs infers commit and compare URLs from git remote
func (g *Generator) inferURLs() error {
	repoURL := g.conf.Repository.URL

	if repoURL == "" {
		var err error
		repoURL, err = g.git.GetRemoteURL()
		if err != nil {
			return err
		}
	}

	if repoURL == "" {
		return nil // no remote, links will be empty
	}

	hostType := git.DetectHostType(repoURL)

	g.commitURL = g.conf.Repository.CommitURL
	if g.commitURL == "" {
		g.commitURL = git.InferCommitURL(repoURL, hostType)
	}

	g.compareURL = g.conf.Repository.CompareURL
	if g.compareURL == "" {
		g.compareURL = git.InferCompareURL(repoURL, hostType)
	}

	return nil
}

// GenerateAll generates the full changelog for all versions
func (g *Generator) GenerateAll() (*Changelog, error) {
	tags, err := g.git.GetSemverTags()
	if err != nil {
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	cl := &Changelog{
		Header: g.conf.Header,
	}

	head, err := g.git.GetHead()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	if len(tags) == 0 {
		// no tags — generate single "Unreleased" version from all commits
		vc, err := g.generateVersion("Unreleased", "", "", head, time.Now())
		if err != nil {
			return nil, err
		}
		if hasCommits(vc) {
			cl.Versions = append(cl.Versions, *vc)
		}
		return cl, nil
	}

	// check for unreleased commits (after latest tag)
	latestTag := tags[0]
	vc, err := g.generateVersion("Unreleased", latestTag.Name, latestTag.Name, "HEAD", time.Now())
	if err != nil {
		return nil, err
	}
	if hasCommits(vc) {
		cl.Versions = append(cl.Versions, *vc)
	}

	// generate each version
	for i := 0; i < len(tags); i++ {
		tag := tags[i]

		var fromRef string
		if i+1 < len(tags) {
			fromRef = tags[i+1].Name
		}

		date := tag.Date
		if date.IsZero() {
			date, _ = g.git.GetTagDate(tag.Name)
		}

		var prevTag string
		if i+1 < len(tags) {
			prevTag = tags[i+1].Name
		}

		vc, err := g.generateVersion(tag.Name, prevTag, fromRef, tag.Name, date)
		if err != nil {
			return nil, err
		}
		if hasCommits(vc) {
			cl.Versions = append(cl.Versions, *vc)
		}
	}

	return cl, nil
}

// GenerateRange generates the changelog between two refs
func (g *Generator) GenerateRange(from, to string) (*Changelog, error) {
	if to == "" {
		to = "HEAD"
	}

	cl := &Changelog{
		Header: g.conf.Header,
	}

	// try to determine version name
	versionName := to
	if to == "HEAD" {
		// check if there's a tag at HEAD
		head, err := g.git.GetHead()
		if err == nil {
			tags, _ := g.git.GetSemverTags()
			for _, t := range tags {
				if t.Hash == head[:len(t.Hash)] || head[:len(t.Hash)] == t.Hash {
					versionName = t.Name
					break
				}
			}
			if versionName == "HEAD" {
				versionName = "Unreleased"
			}
		}
	}

	date := time.Now()
	if to != "HEAD" {
		d, err := g.git.GetTagDate(to)
		if err == nil && !d.IsZero() {
			date = d
		}
	}

	vc, err := g.generateVersion(versionName, from, from, to, date)
	if err != nil {
		return nil, err
	}
	if hasCommits(vc) {
		cl.Versions = append(cl.Versions, *vc)
	}

	return cl, nil
}

// GenerateSmart implements the smart generation logic:
// - If explicit from/to given, use them
// - Otherwise generate full changelog for all versions
func (g *Generator) GenerateSmart(from, to string) (*Changelog, error) {
	// if explicit range given, use it
	if from != "" || to != "" {
		return g.GenerateRange(from, to)
	}

	// generate full changelog for all versions
	return g.GenerateAll()
}

// generateVersion generates a VersionChangelog for a single version
func (g *Generator) generateVersion(version, prevVersion, fromRef, toRef string, date time.Time) (*VersionChangelog, error) {
	commits, err := g.git.GetCommits(fromRef, toRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits for %s: %w", version, err)
	}

	vc := &VersionChangelog{
		Version:    version,
		Date:       date.Format("2006-01-02"),
		CompareURL: git.BuildCompareURL(g.compareURL, prevVersion, version),
		FromRef:    fromRef,
		ToRef:      toRef,
	}

	// if no previous version for compare, and it's the first version
	if prevVersion == "" && len(commits) > 0 {
		firstCommit := commits[len(commits)-1]
		vc.CompareURL = git.BuildCompareURL(g.compareURL, firstCommit.Hash, version)
	}

	// build type lookup
	typeIndex := g.buildTypeIndex()

	// group map: type -> []CommitInfo
	groups := make(map[string][]CommitInfo)
	var breaking []CommitInfo
	var other []CommitInfo

	for _, raw := range commits {
		ci := g.processCommit(raw)

		if ci.Commit == nil {
			// non-conventional commit
			if g.conf.IncludeOther {
				other = append(other, ci)
			}
			continue
		}

		// collect breaking changes
		if g.conf.IncludeBreaking && ci.Commit.IsBreakingChange() {
			breaking = append(breaking, ci)
		}

		commitType := ci.Commit.Type()
		groups[commitType] = append(groups[commitType], ci)
	}

	// build ordered type groups following config order
	for _, tc := range g.conf.Types {
		if tc.Hidden {
			continue
		}
		commits, ok := groups[tc.Type]
		if !ok || len(commits) == 0 {
			continue
		}
		vc.Groups = append(vc.Groups, TypeGroup{
			Type:    tc.Type,
			Header:  tc.Header,
			Commits: commits,
		})
		delete(groups, tc.Type)
	}

	// append any remaining unknown types
	for typ, commits := range groups {
		if len(commits) == 0 {
			continue
		}
		// check if it's hidden
		if tc, ok := typeIndex[typ]; ok && tc.Hidden {
			continue
		}
		vc.Groups = append(vc.Groups, TypeGroup{
			Type:    typ,
			Header:  capitalizeFirst(typ),
			Commits: commits,
		})
	}

	vc.Breaking = breaking
	vc.Other = other

	return vc, nil
}

// processCommit converts a raw git commit into a CommitInfo
func (g *Generator) processCommit(raw git.CommitRaw) CommitInfo {
	ci := CommitInfo{
		Hash:      raw.Hash,
		ShortHash: raw.ShortHash,
		CommitURL: git.BuildCommitURL(g.commitURL, raw.Hash),
		Author:    raw.Author,
		Date:      raw.Date,
	}

	// reconstruct full message for parser
	msg := raw.Subject
	if raw.Body != "" {
		msg = raw.Subject + "\n\n" + raw.Body
	}

	commit, err := g.parser.Parse(msg)
	if err == nil {
		ci.Commit = commit
	}

	// extract issue references
	ci.References = git.ExtractReferences(msg, g.conf.IssuePrefixes)

	return ci
}

// buildTypeIndex creates a map from type name to TypeConfig
func (g *Generator) buildTypeIndex() map[string]TypeConfig {
	index := make(map[string]TypeConfig, len(g.conf.Types))
	for _, tc := range g.conf.Types {
		index[tc.Type] = tc
	}
	return index
}

// hasCommits checks if a VersionChangelog has any commits
func hasCommits(vc *VersionChangelog) bool {
	if vc == nil {
		return false
	}
	for _, g := range vc.Groups {
		if len(g.Commits) > 0 {
			return true
		}
	}
	return len(vc.Breaking) > 0 || len(vc.Other) > 0
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// DebugInfo holds debug information about the repository
type DebugInfo struct {
	Version    string
	Build      string
	GitDir     string
	RemoteURL  string
	HostType   git.HostType
	CommitURL  string
	CompareURL string
	LatestTag  string
	TotalTags  int
	Head       string
	ConfigPath string
	ConfigType string
	Formatter  string
	Types      []TypeConfig
}

// Debug returns debug information about the generator setup
func (g *Generator) Debug() (*DebugInfo, error) {
	head, _ := g.git.GetHead()
	remoteURL, _ := g.git.GetRemoteURL()
	tags, _ := g.git.GetSemverTags()
	latestTag, _ := g.git.GetLatestTag()

	hostType := git.DetectHostType(remoteURL)

	info := &DebugInfo{
		RemoteURL:  remoteURL,
		HostType:   hostType,
		CommitURL:  g.commitURL,
		CompareURL: g.compareURL,
		LatestTag:  latestTag.Name,
		TotalTags:  len(tags),
		Head:       head,
		Formatter:  g.conf.Formatter,
		Types:      g.conf.Types,
	}

	return info, nil
}

package git

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

var (
	errGitNotFound = errors.New("git is required and must be in PATH")
	errNotGitRepo  = errors.New("not a git repository (or any of the parent directories)")
)

// CommitRaw holds raw commit data from git log
type CommitRaw struct {
	Hash      string
	ShortHash string
	Author    string
	Date      time.Time
	Subject   string
	Body      string
}

// TagInfo holds tag data
type TagInfo struct {
	Name string
	Hash string
	Date time.Time
}

// Client wraps git operations via os/exec
type Client struct {
	repoDir string
}

// NewClient creates a new git client for the given repository directory
func NewClient(repoDir string) (*Client, error) {
	_, err := exec.LookPath("git")
	if err != nil {
		return nil, errGitNotFound
	}

	g := &Client{repoDir: repoDir}

	// verify it's a git repo
	_, err = g.run("rev-parse", "--git-dir")
	if err != nil {
		return nil, errNotGitRepo
	}

	return g, nil
}

// run executes a git command and returns trimmed stdout
func (g *Client) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.repoDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return "", fmt.Errorf("git %s: %s", args[0], errMsg)
		}
		return "", fmt.Errorf("git %s: %w", args[0], err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

// GetHead returns the HEAD commit hash
func (g *Client) GetHead() (string, error) {
	return g.run("rev-parse", "HEAD")
}

// GetRemoteURL returns the origin remote URL
func (g *Client) GetRemoteURL() (string, error) {
	out, err := g.run("remote", "get-url", "origin")
	if err != nil {
		return "", nil // no remote is not an error
	}
	return normalizeRemoteURL(out), nil
}

// GetSemverTags returns all semver tags sorted descending (newest first)
func (g *Client) GetSemverTags() ([]TagInfo, error) {
	out, err := g.run("tag", "--list", "--format=%(refname:short)\t%(objectname:short)\t%(*objectname:short)\t%(creatordate:iso-strict)")
	if err != nil {
		return nil, err
	}

	if out == "" {
		return nil, nil
	}

	var tags []TagInfo
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\t", 4)
		if len(parts) < 4 {
			continue
		}

		name := parts[0]

		// ensure it's a valid semver tag
		ver := name
		if !strings.HasPrefix(ver, "v") {
			ver = "v" + ver
		}
		if !semver.IsValid(ver) {
			continue
		}

		// for annotated tags, use the dereferenced commit hash
		hash := parts[1]
		if parts[2] != "" {
			hash = parts[2]
		}

		var t time.Time
		if parts[3] != "" {
			t, _ = time.Parse(time.RFC3339, parts[3])
		}

		tags = append(tags, TagInfo{
			Name: name,
			Hash: hash,
			Date: t,
		})
	}

	// sort descending by semver
	sort.Slice(tags, func(i, j int) bool {
		vi := tags[i].Name
		vj := tags[j].Name
		if !strings.HasPrefix(vi, "v") {
			vi = "v" + vi
		}
		if !strings.HasPrefix(vj, "v") {
			vj = "v" + vj
		}
		return semver.Compare(vi, vj) > 0
	})

	return tags, nil
}

// GetLatestTag returns the latest semver tag, or empty if none
func (g *Client) GetLatestTag() (TagInfo, error) {
	tags, err := g.GetSemverTags()
	if err != nil {
		return TagInfo{}, err
	}
	if len(tags) == 0 {
		return TagInfo{}, nil
	}
	return tags[0], nil
}

// GetFirstCommit returns the first commit hash in the repository
func (g *Client) GetFirstCommit() (string, error) {
	return g.run("rev-list", "--max-parents=0", "HEAD")
}

const (
	commitSep = "---COMMIT_SEP---"
	fieldSep  = "---FIELD_SEP---"
)

// commitLogFormat is the git log format string
var commitLogFormat = strings.Join([]string{
	commitSep,
	"%H" + fieldSep + "%h" + fieldSep + "%an" + fieldSep + "%aI" + fieldSep + "%s" + fieldSep + "%b",
}, "")

// GetCommits returns commits between two refs
// fromRef is exclusive, toRef is inclusive
// if fromRef is empty, returns all commits up to toRef
// if toRef is empty, defaults to HEAD
func (g *Client) GetCommits(fromRef, toRef string) ([]CommitRaw, error) {
	if toRef == "" {
		toRef = "HEAD"
	}

	var refRange string
	if fromRef == "" {
		refRange = toRef
	} else {
		refRange = fromRef + ".." + toRef
	}

	out, err := g.run("log", "--format="+commitLogFormat, "--no-merges", refRange)
	if err != nil {
		return nil, err
	}

	if out == "" {
		return nil, nil
	}

	return parseGitLog(out), nil
}

// GetTagDate returns the date of a tag
func (g *Client) GetTagDate(tag string) (time.Time, error) {
	out, err := g.run("log", "-1", "--format=%aI", tag)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.Parse(time.RFC3339, strings.TrimSpace(out))
	if err != nil {
		return time.Time{}, nil
	}
	return t, nil
}

// HasChangelog checks if CHANGELOG.md exists
func (g *Client) HasChangelog() bool {
	_, err := g.run("ls-files", "CHANGELOG.md")
	return err == nil
}

// parseGitLog parses the output of git log
func parseGitLog(output string) []CommitRaw {
	chunks := strings.Split(output, commitSep)
	var commits []CommitRaw

	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		parts := strings.SplitN(chunk, fieldSep, 6)
		if len(parts) < 5 {
			continue
		}

		hash := strings.TrimSpace(parts[0])
		shortHash := strings.TrimSpace(parts[1])
		author := strings.TrimSpace(parts[2])
		dateStr := strings.TrimSpace(parts[3])
		subject := strings.TrimSpace(parts[4])

		var body string
		if len(parts) > 5 {
			body = strings.TrimSpace(parts[5])
		}

		var date time.Time
		if dateStr != "" {
			date, _ = time.Parse(time.RFC3339, dateStr)
		}

		commits = append(commits, CommitRaw{
			Hash:      hash,
			ShortHash: shortHash,
			Author:    author,
			Date:      date,
			Subject:   subject,
			Body:      body,
		})
	}

	return commits
}

// ssh pattern: git@host:user/repo.git
var sshURLRegexp = regexp.MustCompile(`^[\w-]+@([\w.-]+):([\w./-]+?)(?:\.git)?$`)

// normalizeRemoteURL converts git remote URL to HTTPS URL
func normalizeRemoteURL(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)

	// handle SSH URLs: git@github.com:user/repo.git
	if matches := sshURLRegexp.FindStringSubmatch(rawURL); matches != nil {
		return "https://" + matches[1] + "/" + matches[2]
	}

	// handle HTTPS URLs: remove .git suffix
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	parsed.Path = strings.TrimSuffix(parsed.Path, ".git")

	// ensure https
	if parsed.Scheme == "http" {
		parsed.Scheme = "https"
	}

	// remove any userinfo
	parsed.User = nil

	return parsed.String()
}

// HostType represents the type of git hosting service
type HostType string

const (
	HostGitHub    HostType = "github"
	HostGitLab    HostType = "gitlab"
	HostBitbucket HostType = "bitbucket"
	HostAzure     HostType = "azure"
	HostGeneric   HostType = "generic"
)

// DetectHostType detects the hosting service from a URL
func DetectHostType(repoURL string) HostType {
	lower := strings.ToLower(repoURL)

	switch {
	case strings.Contains(lower, "github.com"):
		return HostGitHub
	case strings.Contains(lower, "gitlab.com") || strings.Contains(lower, "gitlab"):
		return HostGitLab
	case strings.Contains(lower, "bitbucket.org") || strings.Contains(lower, "bitbucket"):
		return HostBitbucket
	case strings.Contains(lower, "dev.azure.com") || strings.Contains(lower, "visualstudio.com"):
		return HostAzure
	default:
		return HostGeneric
	}
}

// InferCommitURL returns the commit URL template for the given host
func InferCommitURL(repoURL string, hostType HostType) string {
	switch hostType {
	case HostGitLab:
		return repoURL + "/-/commit/{{hash}}"
	case HostBitbucket:
		return repoURL + "/commits/{{hash}}"
	case HostAzure:
		return repoURL + "#/commit/{{hash}}"
	default: // GitHub, generic
		return repoURL + "/commit/{{hash}}"
	}
}

// InferCompareURL returns the compare URL template for the given host
func InferCompareURL(repoURL string, hostType HostType) string {
	switch hostType {
	case HostGitLab:
		return repoURL + "/-/compare/{{from}}...{{to}}"
	case HostBitbucket:
		return repoURL + "/compare/{{from}}..{{to}}"
	case HostAzure:
		return repoURL + "#/compare?head=true&sourceBranch={{to}}&targetBranch={{from}}"
	default: // GitHub, generic
		return repoURL + "/compare/{{from}}...{{to}}"
	}
}

// BuildCommitURL replaces templates in commit URL
func BuildCommitURL(tmpl, hash string) string {
	if tmpl == "" {
		return ""
	}
	return strings.ReplaceAll(tmpl, "{{hash}}", hash)
}

// BuildCompareURL replaces templates in compare URL
func BuildCompareURL(tmpl, from, to string) string {
	if tmpl == "" {
		return ""
	}
	r := strings.ReplaceAll(tmpl, "{{from}}", from)
	r = strings.ReplaceAll(r, "{{to}}", to)
	return r
}

// ExtractReferences extracts issue references from a commit message
func ExtractReferences(msg string, prefixes []string) []string {
	if len(prefixes) == 0 {
		return nil
	}

	var refs []string
	seen := make(map[string]struct{})

	for _, prefix := range prefixes {
		pattern := regexp.QuoteMeta(prefix) + `(\d+)`
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		matches := re.FindAllStringSubmatch(msg, -1)
		for _, m := range matches {
			ref := prefix + m[1]
			if _, ok := seen[ref]; !ok {
				refs = append(refs, ref)
				seen[ref] = struct{}{}
			}
		}
	}

	return refs
}

#!/usr/bin/env bash

# integration_test.sh - real-world git integration tests for commitlint
#
# Usage:
#   ./integration_test.sh /path/to/commitlint [/path/to/tmpdir]
#
# If tmpdir is omitted, a temporary directory is created automatically.
# The script creates git repos, makes various commits, and verifies
# commitlint output for each scenario.

set -euo pipefail

# ─── Arguments ───────────────────────────────────────────────────────────────

COMMITLINT="${1:?Usage: $0 /path/to/commitlint [tmpdir]}"
TMPDIR_ROOT="${2:-$(mktemp -d)}"

if [[ ! -x "$COMMITLINT" ]]; then
  echo "ERROR: $COMMITLINT is not executable" >&2
  exit 1
fi

COMMITLINT="$(realpath "$COMMITLINT")"

# ─── Counters ────────────────────────────────────────────────────────────────

PASS=0
FAIL=0
TOTAL=0

pass() {
  PASS=$((PASS + 1))
  TOTAL=$((TOTAL + 1))
  echo "  ✔ $1"
}

fail() {
  FAIL=$((FAIL + 1))
  TOTAL=$((TOTAL + 1))
  echo "  ✘ $1"
}

# Run commitlint lint on a message string. Sets $EXIT_CODE and $OUTPUT.
run_lint() {
  local msg="$1"
  set +e
  OUTPUT=$(echo "$msg" | "$COMMITLINT" lint 2>&1)
  EXIT_CODE=$?
  set -e
}

# Expect commitlint to succeed (exit 0)
expect_pass() {
  local desc="$1" msg="$2"
  run_lint "$msg"
  if [[ $EXIT_CODE -eq 0 ]]; then
    pass "$desc"
  else
    fail "$desc (expected pass, got exit $EXIT_CODE)"
    echo "       msg: $msg"
    echo "       out: $OUTPUT"
  fi
}

# Expect commitlint to fail (exit != 0)
expect_fail() {
  local desc="$1" msg="$2"
  run_lint "$msg"
  if [[ $EXIT_CODE -ne 0 ]]; then
    pass "$desc"
  else
    fail "$desc (expected fail, got exit 0)"
    echo "       msg: $msg"
    echo "       out: $OUTPUT"
  fi
}

# ─── Setup temp repo ────────────────────────────────────────────────────────

REPO="$TMPDIR_ROOT/test-repo"
mkdir -p "$REPO"
cd "$REPO"
git init -q
git config user.email "test@test.com"
git config user.name "Test User"
# Disable any global commit hooks so they don't interfere
git config core.hooksPath /dev/null

# Create initial file so we have something to commit
echo "init" > file.txt
git add file.txt
git commit -q -m "feat: initial project setup with file"

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  commitlint integration tests"
echo "  binary:  $COMMITLINT"
echo "  tmpdir:  $TMPDIR_ROOT"
echo "═══════════════════════════════════════════════════════════"

# ═════════════════════════════════════════════════════════════════════════════
# 1. Valid conventional commits
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 1. Valid conventional commits ──"

expect_pass "feat with description" \
  "feat: add new login page"

expect_pass "fix with scope" \
  "fix(auth): resolve token expiry issue"

expect_pass "docs change" \
  "docs: update README with examples"

expect_pass "chore with long description" \
  "chore: update dependencies to latest"

expect_pass "refactor with scope" \
  "refactor(core): simplify parser logic"

expect_pass "ci change" \
  "ci: add GitHub Actions workflow"

expect_pass "build change" \
  "build: upgrade Go to 1.24"

expect_pass "test addition" \
  "test: add unit tests for parser"

expect_pass "perf improvement" \
  "perf: optimize database queries"

expect_pass "style change" \
  "style: format code with gofmt"

expect_pass "revert type" \
  "revert: undo accidental deletion"

expect_pass "multiline body" \
  "feat: add user profiles

This adds a new user profile page with avatar support,
bio editing, and social links."

# ═════════════════════════════════════════════════════════════════════════════
# 2. Invalid conventional commits
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 2. Invalid conventional commits ──"

expect_fail "unknown type" \
  "fear: this type does not exist"

expect_fail "missing colon" \
  "feat add something without colon"

expect_fail "too short header" \
  "fix: x"

expect_fail "random text" \
  "just some random text without structure"

expect_fail "lowercase type with bad format" \
  "FEAT: uppercase type"

# ═════════════════════════════════════════════════════════════════════════════
# 3. Ignored: merge commits
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 3. Ignored: merge commits ──"

expect_pass "GitHub PR merge" \
  "Merge pull request #123 from owner/feature-branch"

expect_pass "merge branch" \
  "Merge branch 'feature-x'"

expect_pass "merge branch into" \
  "Merge branch 'release/1.0' into main"

expect_pass "merge tag" \
  "Merge tag 'v1.0.0'"

expect_pass "merge remote-tracking" \
  "Merge remote-tracking branch 'origin/main'"

expect_pass "merge X into Y" \
  "Merge feature into main"

expect_pass "Azure DevOps PR" \
  "Merged PR #456: Add new feature"

expect_pass "Azure DevOps PR no hash" \
  "Merged PR 789: Fix bug"

expect_pass "Merged X into Y" \
  "Merged feature-branch into main"

expect_pass "Merged X in Y" \
  "Merged feature-branch in main"

expect_pass "multiline merge" \
  "Merge pull request #100 from org/branch

Additional context about the merge"

# ═════════════════════════════════════════════════════════════════════════════
# 4. Ignored: revert / reapply
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 4. Ignored: revert / reapply ──"

expect_pass "Revert with quotes" \
  'Revert "feat: add new feature"'

expect_pass "revert lowercase" \
  'revert "fix: something"'

expect_pass "Revert plain" \
  "Revert some commit"

expect_pass "Reapply with quotes" \
  'Reapply "feat: add feature"'

expect_pass "reapply lowercase" \
  'reapply "fix: something"'

# ═════════════════════════════════════════════════════════════════════════════
# 5. Ignored: fixup / squash / amend
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 5. Ignored: fixup / squash / amend ──"

expect_pass "fixup!" \
  "fixup! feat: add something"

expect_pass "squash!" \
  "squash! fix: repair something"

expect_pass "amend!" \
  "amend! chore: update deps"

# ═════════════════════════════════════════════════════════════════════════════
# 6. Ignored: automatic merges / initial commit
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 6. Ignored: automatic merges / initial commit ──"

expect_pass "Automatic merge" \
  "Automatic merge from release/1.0 to main"

expect_pass "Auto-merged" \
  "Auto-merged feature-x into main"

expect_pass "Initial commit" \
  "Initial commit"

# ═════════════════════════════════════════════════════════════════════════════
# 7. NOT ignored (look-alikes)
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 7. NOT ignored (look-alikes) ──"

expect_fail "lowercase merge" \
  "merge something random"

expect_fail "fixup without bang" \
  "fixup something without bang"

expect_fail "lowercase initial commit" \
  "initial commit"

expect_fail "Initial commit with suffix" \
  "Initial commit with extra text"

# ═════════════════════════════════════════════════════════════════════════════
# 8. Real git merge commit (actual git operations)
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 8. Real git merge commit ──"

# Create a branch, make a commit, merge it
cd "$REPO"
git checkout -q -b test-feature
echo "feature" > feature.txt
git add feature.txt
git commit -q -m "feat: add feature file"

git checkout -q master 2>/dev/null || git checkout -q main
# Do a no-ff merge to generate a real merge commit message
git merge --no-ff test-feature -q -m "Merge branch 'test-feature'"

# Get the merge commit message
MERGE_MSG=$(git log -1 --format=%B)
expect_pass "real git merge commit" "$MERGE_MSG"

# ═════════════════════════════════════════════════════════════════════════════
# 9. Config file tests
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 9. Config file tests ──"

# Test with disable-default-ignores
CONF_DIR="$TMPDIR_ROOT/conf-tests"
mkdir -p "$CONF_DIR"

cat > "$CONF_DIR/strict.yaml" << 'YAML'
min-version: v0.9.0
formatter: default
rules:
  - header-min-length
  - header-max-length
  - type-enum
severity:
  default: error
settings:
  header-min-length:
    argument: 10
  header-max-length:
    argument: 50
  type-enum:
    argument:
      - feat
      - fix
      - docs
      - chore
disable-default-ignores: true
YAML

# Merge commit should FAIL with defaults disabled
set +e
OUTPUT=$(echo "Merge pull request #123 from owner/branch" | "$COMMITLINT" lint --config="$CONF_DIR/strict.yaml" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -ne 0 ]]; then
  pass "merge fails with disable-default-ignores: true"
else
  fail "merge should fail with disable-default-ignores: true"
fi

# Test with custom ignore pattern
cat > "$CONF_DIR/custom.yaml" << 'YAML'
min-version: v0.9.0
formatter: default
rules:
  - header-min-length
  - header-max-length
  - type-enum
severity:
  default: error
settings:
  header-min-length:
    argument: 10
  header-max-length:
    argument: 50
  type-enum:
    argument:
      - feat
      - fix
      - docs
      - chore
ignores:
  - "^WIP "
  - "^TICKET-\\d+"
YAML

# WIP should pass (custom pattern)
set +e
OUTPUT=$(echo "WIP adding new feature" | "$COMMITLINT" lint --config="$CONF_DIR/custom.yaml" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -eq 0 ]]; then
  pass "WIP passes with custom ignores"
else
  fail "WIP should pass with custom ignores"
fi

# Merge should ALSO pass (defaults still active, additive)
set +e
OUTPUT=$(echo "Merge branch 'feature'" | "$COMMITLINT" lint --config="$CONF_DIR/custom.yaml" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -eq 0 ]]; then
  pass "merge still passes with additive custom ignores"
else
  fail "merge should still pass (custom ignores are additive to defaults)"
fi

# TICKET pattern
set +e
OUTPUT=$(echo "TICKET-42 implement dashboard" | "$COMMITLINT" lint --config="$CONF_DIR/custom.yaml" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -eq 0 ]]; then
  pass "TICKET-42 passes with custom ignores"
else
  fail "TICKET-42 should pass with custom ignores"
fi

# Config check should succeed
set +e
OUTPUT=$("$COMMITLINT" config check "$CONF_DIR/custom.yaml" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -eq 0 ]]; then
  pass "config check passes for valid config"
else
  fail "config check should pass for valid config"
fi

# ═════════════════════════════════════════════════════════════════════════════
# 10. Pipe / file input modes
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "── 10. Input modes ──"

# Pipe mode (already tested above, but explicit)
set +e
OUTPUT=$(echo "feat: add something new here" | "$COMMITLINT" lint 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -eq 0 ]]; then
  pass "pipe input mode"
else
  fail "pipe input mode"
fi

# File input mode
MSG_FILE="$TMPDIR_ROOT/commit_msg.txt"
echo "feat: add something via file" > "$MSG_FILE"
set +e
OUTPUT=$("$COMMITLINT" lint --message="$MSG_FILE" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -eq 0 ]]; then
  pass "file input mode (--message)"
else
  fail "file input mode (--message)"
fi

# Redirect input mode
set +e
OUTPUT=$("$COMMITLINT" lint < "$MSG_FILE" 2>&1)
EXIT_CODE=$?
set -e
if [[ $EXIT_CODE -eq 0 ]]; then
  pass "redirect input mode"
else
  fail "redirect input mode"
fi

# ═════════════════════════════════════════════════════════════════════════════
# Summary
# ═════════════════════════════════════════════════════════════════════════════

echo ""
echo "═══════════════════════════════════════════════════════════"
echo "  Results: $PASS passed, $FAIL failed ($TOTAL total)"
echo "═══════════════════════════════════════════════════════════"

# Cleanup
rm -rf "$TMPDIR_ROOT"

if [[ $FAIL -gt 0 ]]; then
  exit 1
fi
exit 0

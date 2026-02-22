[![PkgGoDev](https://pkg.go.dev/badge/github.com/conventionalcommit/commitlint)](https://pkg.go.dev/github.com/conventionalcommit/commitlint)

# commitlint

commitlint checks if your commit message meets the [conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/)

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

- [Why Use Conventional Commits?](https://www.conventionalcommits.org/en/v1.0.0/#why-use-conventional-commits)

### Table of Contents

- [commitlint](#commitlint)
    - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Releases](#releases)
    - [Using go](#using-go)
  - [Setup](#setup)
    - [Manual](#manual)
    - [Remove](#remove)
  - [Quick Test](#quick-test)
  - [Commands](#commands)
    - [config](#config)
    - [lint](#lint)
      - [Config Precedence](#config-precedence)
      - [Message Precedence](#message-precedence)
    - [hook](#hook)
    - [debug](#debug)
  - [Default Config](#default-config)
    - [Commit Types](#commit-types)
  - [Ignore Patterns](#ignore-patterns)
    - [Default Ignore Patterns](#default-ignore-patterns)
    - [Custom Ignore Patterns](#custom-ignore-patterns)
    - [Disabling Default Ignores](#disabling-default-ignores)
  - [Available Rules](#available-rules)
    - [Length rules](#length-rules)
    - [Enum / allow-list rules](#enum--allow-list-rules)
    - [Charset rules](#charset-rules)
    - [Case rules](#case-rules)
    - [Empty / presence rules](#empty--presence-rules)
    - [Full-stop rules](#full-stop-rules)
    - [Leading-blank rules](#leading-blank-rules)
    - [Header formatting rules](#header-formatting-rules)
    - [Trailer / sign-off rules](#trailer--sign-off-rules)
    - [Breaking change rules](#breaking-change-rules)
  - [Available Formatters](#available-formatters)
  - [Programmatic Usage](#programmatic-usage)
    - [One-liner with default config](#one-liner-with-default-config)
    - [Full control with default config](#full-control-with-default-config)
    - [Lint with a config file](#lint-with-a-config-file)
    - [Custom rules](#custom-rules)
    - [Custom formatters](#custom-formatters)
  - [FAQ](#faq)
  - [License](#license)

## Installation

### Releases

Download binary from [releases](https://github.com/conventionalcommit/commitlint/releases) and add it to your `PATH`

### Using go

```bash
go install github.com/conventionalcommit/commitlint@latest
```

## Setup

- Enable for a single git repository, `cd` to repository directory

```bash
commitlint init
```

- Enable globally for all git repositories

```bash
commitlint init --global
```

- to customize hooks destination pass `--hookspath` with desired location

```bash
commitlint init --hookspath /path/to/hooks
commitlint init --global --hookspath /path/to/hooks
```

### Manual

- run `commitlint hook` to create `.commitlint/hooks` containing git hooks
  - pass `--hookspath` or `-p` to customize the hooks output path
- To enable in single repo
  - run `git config core.hooksPath /path/to/.commitlint/hooks`
- To enable globally
  - run `git config --global core.hooksPath /path/to/.commitlint/hooks`

### Remove

- To remove hooks from a single repository

```bash
commitlint remove
```

- To remove hooks globally

```bash
commitlint remove --global
```

Both commands ask for confirmation before unsetting `core.hooksPath` in git config. Hook files are left intact.

## Quick Test

- Valid commit message

```bash
echo "feat: good commit message" | commitlint lint
# ✔ commit message
```

- Invalid commit message

```bash
echo "fear: do not fear for commit message" | commitlint lint
#   ❌ type-enum: type 'fear' is not allowed, you can use one of [build chore ci docs feat fix merge perf refactor revert style test]
```

## Commands

### config

- To create a config file, run `commitlint config create`, this will create `.commitlint.yaml` with only the enabled rules and their settings (compact format)

- To create a config file with **all** rules and settings written out (including disabled ones), run `commitlint config create --all`

- To validate a config file, run `commitlint config check /path/to/conf.yaml`

### lint

To lint a message, you can use any one of the following
- run `commitlint lint --message=file`
- run `echo "message" | commitlint lint`
- run `commitlint lint < file`

`commitlint lint` follows below order for `config` and `message`

#### Config Precedence

- config file passed to `--config` command-line argument
- `COMMITLINT_CONFIG` env variable
- config file in current directory or git repo root in the below order
  - .commitlint.yml
  - .commitlint.yaml
  - commitlint.yml
  - commitlint.yaml
- [default config](#default-config)

#### Message Precedence

- `stdin` pipe stream
- commit message file passed to `--message` command-line argument
- `.git/COMMIT_EDITMSG` in current directory

### hook

- To create hook files, run `commitlint hook`
  - pass `--hookspath` or `-p` to customize the hooks output directory
  - pass `--replace` or `-r` to overwrite existing hook files

### debug

  To prints useful information for debugging commitlint

  run `commitlint debug`

## Default Config

```yaml
min-version: v0.11.0
formatter: default
rules:
- header-min-length
- header-max-length
- body-max-line-length
- footer-max-line-length
- type-enum
severity:
  default: error
  rules: {}
settings:
  body-max-line-length:
    argument: 100
    flags: {}
  footer-max-line-length:
    argument: 100
    flags: {}
  header-max-length:
    argument: 72
    flags: {}
  header-min-length:
    argument: 10
    flags: {}
  type-enum:
    argument:
    - feat
    - fix
    - docs
    - style
    - refactor
    - perf
    - test
    - build
    - ci
    - chore
    - revert
    flags: {}
disable-default-ignores: false
ignores: []
```

### Commit Types

Commonly used commit types

| Type     | Description                                                                      |
|:---------|:---------------------------------------------------------------------------------|
| feat     | A new feature                                                                    |
| fix      | A bug fix                                                                        |
| docs     | Documentation only changes                                                       |
| style    | Changes that do not affect the meaning of the code (white-space, formatting etc) |
| refactor | A code change that neither fixes a bug nor adds a feature                        |
| perf     | A code change that improves performance                                          |
| test     | Adding missing tests or correcting existing tests                                |
| build    | Changes that affect the build system or external dependencies                    |
| ci       | Changes to our CI configuration files and scripts                                |
| chore    | Other changes that don't modify src or test files                                |
| revert   | Reverts a previous commit                                                        |

## Ignore Patterns

commitlint automatically skips linting for commit messages generated by git (merges, reverts, fixups, etc.).
If the **first line** of a commit message matches any ignore pattern, linting is skipped entirely.

### Default Ignore Patterns

The following patterns are enabled by default
(source: [`config/default.go`](config/default.go)):

| Pattern                              | Matches                               |
|:-------------------------------------|:--------------------------------------|
| `^Merge pull request #\d+`           | GitHub pull request merges            |
| `^Merge .+ into .+`                  | Generic merge (X into Y)              |
| `^Merge branch '.+'`                 | `git merge` branch                    |
| `^Merge tag '.+'`                    | `git merge` tag                       |
| `^Merge remote-tracking branch '.+'` | `git merge` remote-tracking branch    |
| `^Merged .+ (in\|into) .+`           | Azure DevOps / Bitbucket merged       |
| `^Merged PR #?\d+`                   | Azure DevOps pull request             |
| `^(R\|r)evert `                      | `git revert`                          |
| `^(R\|r)eapply `                     | `git reapply`                         |
| `^(amend\|fixup\|squash)! `          | `git commit --fixup/--squash/--amend` |
| `^Automatic merge`                   | Automatic merges                      |
| `^Auto-merged .+ into .+`            | Auto-merged branches                  |
| `^Initial commit$`                   | Initial commit (exact match)          |

### Custom Ignore Patterns

Add your own patterns in the config file under `ignores:`. User-defined patterns are
**additive**, they are checked alongside the built-in defaults.

```yaml
ignores:
  - "^WIP "
  - "^TICKET-\\d+"
```

### Disabling Default Ignores

If you want **only** your custom patterns (no built-in defaults), set `disable-default-ignores: true`:

```yaml
disable-default-ignores: true
ignores:
  - "^WIP "
```

## Available Rules

Rules marked **✅ enabled** are active by default. All others can be opted into via the `rules:` list in your config.

### Length rules

| name                     | argument | flags | description                       | default         |
|:-------------------------|:---------|:------|:----------------------------------|:----------------|
| `header-min-length`      | int      | n/a   | min length of header (first line) | ✅ enabled (10)  |
| `header-max-length`      | int      | n/a   | max length of header (first line) | ✅ enabled (72)  |
| `body-min-length`        | int      | n/a   | min length of body                | N/A             |
| `body-max-length`        | int      | n/a   | max length of body                | N/A             |
| `body-max-line-length`   | int      | n/a   | max length of each line in body   | ✅ enabled (100) |
| `footer-min-length`      | int      | n/a   | min length of footer              | N/A             |
| `footer-max-length`      | int      | n/a   | max length of footer              | N/A             |
| `footer-max-line-length` | int      | n/a   | max length of each line in footer | ✅ enabled (100) |
| `type-min-length`        | int      | n/a   | min length of type                | N/A             |
| `type-max-length`        | int      | n/a   | max length of type                | N/A             |
| `scope-min-length`       | int      | n/a   | min length of scope               | N/A             |
| `scope-max-length`       | int      | n/a   | max length of scope               | N/A             |
| `description-min-length` | int      | n/a   | min length of description         | N/A             |
| `description-max-length` | int      | n/a   | max length of description         | N/A             |

### Enum / allow-list rules

| name               | argument                   | flags               | description                             | default   |
|:-------------------|:---------------------------|:--------------------|:----------------------------------------|:----------|
| `type-enum`        | `[]string`                 | n/a                 | restrict type to given list of strings  | ✅ enabled |
| `scope-enum`       | `[]string`                 | `allow-empty: bool` | restrict scope to given list of strings | N/A       |
| `footer-enum`      | `[]string`                 | n/a                 | restrict footer token to given list     | N/A       |
| `footer-type-enum` | `[]{token, types, values}` | n/a                 | enforce footer notes for given type     | N/A       |

### Charset rules

| name            | argument | flags | description                     | default |
|:----------------|:---------|:------|:--------------------------------|:--------|
| `type-charset`  | string   | n/a   | restrict type to given charset  | N/A     |
| `scope-charset` | string   | n/a   | restrict scope to given charset | N/A     |

### Case rules

All case rules accept one of: `lower-case`, `upper-case`, `camel-case`, `kebab-case`, `pascal-case`, `sentence-case`, `snake-case`, `start-case`.

| name               | argument | flags | description                                | default |
|:-------------------|:---------|:------|:-------------------------------------------|:--------|
| `type-case`        | string   | n/a   | enforce case format on type                | N/A     |
| `scope-case`       | string   | n/a   | enforce case format on scope (skips empty) | N/A     |
| `description-case` | string   | n/a   | enforce case format on description         | N/A     |
| `body-case`        | string   | n/a   | enforce case format on entire body         | N/A     |
| `header-case`      | string   | n/a   | enforce case format on full header         | N/A     |

### Empty / presence rules

These rules enforce that a field is **not empty**.

| name                | argument | flags | description                   | default |
|:--------------------|:---------|:------|:------------------------------|:--------|
| `type-empty`        | n/a      | n/a   | type must not be empty        | N/A     |
| `scope-empty`       | n/a      | n/a   | scope must not be empty       | N/A     |
| `body-empty`        | n/a      | n/a   | body must not be empty        | N/A     |
| `footer-empty`      | n/a      | n/a   | footer must not be empty      | N/A     |
| `description-empty` | n/a      | n/a   | description must not be empty | N/A     |

### Full-stop rules

Check that a field does **not** end with a given character (default `"."`).

| name                    | argument | flags | description                                      | default |
|:------------------------|:---------|:------|:-------------------------------------------------|:--------|
| `header-full-stop`      | string   | n/a   | header must not end with given char (e.g. `"."`) | N/A     |
| `body-full-stop`        | string   | n/a   | body must not end with given char                | N/A     |
| `description-full-stop` | string   | n/a   | description must not end with given char         | N/A     |

### Leading-blank rules

Enforce that a blank line separates commit sections (conventional commits spec).

| name                   | argument | flags | description                             | default |
|:-----------------------|:---------|:------|:----------------------------------------|:--------|
| `body-leading-blank`   | n/a      | n/a   | body must be preceded by a blank line   | N/A     |
| `footer-leading-blank` | n/a      | n/a   | footer must be preceded by a blank line | N/A     |

### Header formatting rules

| name          | argument | flags | description                                         | default |
|:--------------|:---------|:------|:----------------------------------------------------|:--------|
| `header-trim` | n/a      | n/a   | header must not have leading or trailing whitespace | N/A     |

### Trailer / sign-off rules

The argument is the trailer token. A trailing `:` is accepted and stripped automatically,
so `"Signed-off-by"` and `"Signed-off-by:"` are equivalent.

| name             | argument | flags | description                                                          | default |
|:-----------------|:---------|:------|:---------------------------------------------------------------------|:--------|
| `signed-off-by`  | string   | n/a   | commit must have a footer note whose token matches (e.g. `"Signed-off-by"`) | N/A     |
| `trailer-exists` | string   | n/a   | commit must have a footer note whose token matches (e.g. `"Co-authored-by"`) | N/A     |

### Breaking change rules

| name                               | argument | flags | description                                                                                                    | default |
|:-----------------------------------|:---------|:------|:---------------------------------------------------------------------------------------------------------------|:--------|
| `breaking-change-exclamation-mark` | n/a      | n/a   | XNOR: either both `!` in header and `BREAKING CHANGE` in footer are present, or neither N/A not just one alone | N/A     |

## Available Formatters

- default

```
commitlint

→ input: "fear: do not fear for ..."

Errors:
  ❌ type-enum: type 'fear' is not allowed, you can use one of [build chore ci docs feat fix perf refactor revert style test]

Total 1 errors, 0 warnings, 0 other severities
```

- JSON

```json
{"input":"fear: do not fear for commit message","issues":[{"description":"type 'fear' is not allowed, you can use one of [build chore ci docs feat fix perf refactor revert style test]","name":"type-enum","severity":"error"}]}
```

## Programmatic Usage

All public packages are importable. The module path is `github.com/conventionalcommit/commitlint`.

```bash
go get github.com/conventionalcommit/commitlint@latest
```

Key packages:

| Package | Purpose |
|:--------|:--------|
| `config` | Parse config files, build a `Linter`, access defaults |
| `lint` | Core types: `Linter`, `Rule`, `Formatter`, `Config`, `Result`, `Issue` |
| `registry` | Register and look up custom rules / formatters |
| `rule` | Built-in rule implementations |
| `formatter` | Built-in formatters (`default`, `json`) |

### One-liner with default config

The simplest entry point — no config file required:

```go
package main

import (
    "fmt"
    "github.com/conventionalcommit/commitlint/config"
)

func main() {
    result, err := config.LintMessage("feat: add login page")
    if err != nil {
        panic(err)
    }

    for _, issue := range result.Issues() {
        fmt.Printf("%s: %s: %s\n", issue.Severity(), issue.RuleName(), issue.Description())
    }

    if len(result.Issues()) == 0 {
        fmt.Println("commit message is valid")
    }
}
```

### Full control with default config

Build the linter yourself for more control (e.g. to swap the formatter):

```go
package main

import (
    "fmt"
    "github.com/conventionalcommit/commitlint/config"
    "github.com/conventionalcommit/commitlint/formatter"
)

func main() {
    conf := config.NewDefault()
    // optionally customise conf here

    linter, err := config.NewLinter(conf)
    if err != nil {
        panic(err)
    }

    result, err := linter.ParseAndLint("feat: add login page")
    if err != nil {
        panic(err)
    }

    out, err := (&formatter.JSONFormatter{}).Format(result)
    if err != nil {
        panic(err)
    }
    fmt.Println(out)
}
```

### Lint with a config file

Load a `.commitlint.yaml` and lint against it:

```go
package main

import (
    "fmt"
    "github.com/conventionalcommit/commitlint/config"
)

func main() {
    conf, err := config.Parse(".commitlint.yaml")
    if err != nil {
        panic(err)
    }

    linter, err := config.NewLinter(conf)
    if err != nil {
        panic(err)
    }

    result, err := linter.ParseAndLint("feat: add login page")
    if err != nil {
        panic(err)
    }

    for _, issue := range result.Issues() {
        fmt.Printf("%s: %s\n", issue.RuleName(), issue.Description())
    }
}
```

### Custom rules

Implement the `lint.Rule` interface and register it before building a linter:

```go
package main

import (
    "fmt"
    "github.com/conventionalcommit/commitlint/config"
    "github.com/conventionalcommit/commitlint/lint"
    "github.com/conventionalcommit/commitlint/registry"
)

// NoWIPRule rejects commit messages whose description starts with "WIP".
type NoWIPRule struct{}

func (r *NoWIPRule) Name() string { return "no-wip" }
func (r *NoWIPRule) Apply(setting lint.RuleSetting) error { return nil }
func (r *NoWIPRule) Validate(commit lint.Commit) (*lint.Issue, error) {
    if len(commit.Description()) >= 3 && commit.Description()[:3] == "WIP" {
        return lint.NewIssue("description must not start with WIP"), nil
    }
    return nil, nil
}

func main() {
    if err := registry.RegisterRule(&NoWIPRule{}); err != nil {
        panic(err)
    }

    conf := config.NewDefault()
    conf.Rules = append(conf.Rules, "no-wip")
    conf.Settings["no-wip"] = lint.RuleSetting{}

    linter, err := config.NewLinter(conf)
    if err != nil {
        panic(err)
    }

    result, err := linter.ParseAndLint("feat: WIP do not merge")
    if err != nil {
        panic(err)
    }

    for _, issue := range result.Issues() {
        fmt.Printf("%s: %s\n", issue.RuleName(), issue.Description())
    }
}
```

### Custom formatters

Implement `lint.Formatter` and register it:

```go
package main

import (
    "fmt"
    "strings"
    "github.com/conventionalcommit/commitlint/config"
    "github.com/conventionalcommit/commitlint/lint"
    "github.com/conventionalcommit/commitlint/registry"
)

type SimpleFormatter struct{}

func (f *SimpleFormatter) Name() string { return "simple" }
func (f *SimpleFormatter) Format(result *lint.Result) (string, error) {
    if len(result.Issues()) == 0 {
        return "ok", nil
    }
    var sb strings.Builder
    for _, issue := range result.Issues() {
        fmt.Fprintf(&sb, "[%s] %s: %s\n", issue.Severity(), issue.RuleName(), issue.Description())
    }
    return sb.String(), nil
}

func main() {
    if err := registry.RegisterFormatter(&SimpleFormatter{}); err != nil {
        panic(err)
    }

    conf := config.NewDefault()
    conf.Formatter = "simple"

    format, err := config.GetFormatter(conf)
    if err != nil {
        panic(err)
    }

    linter, err := config.NewLinter(conf)
    if err != nil {
        panic(err)
    }

    result, err := linter.ParseAndLint("bad message")
    if err != nil {
        panic(err)
    }

    out, err := format.Format(result)
    if err != nil {
        panic(err)
    }
    fmt.Print(out)
}
```

## FAQ

- How to have custom config for each repository?

  Place `.commitlint.yaml` file in repo root directory. linter follows [config precedence](#config-precedence).

  To create a sample config, run `commitlint config create` (or `commitlint config create --all` to include all available settings)

- How can I skip lint check for a commit?

  use `--no-verify` flag with `git commit` which skips commit hooks

- How does commitlint handle merge / revert commits?

  commitlint ships with [built-in ignore patterns](#default-ignore-patterns) that
  automatically skip linting for merge commits, reverts, fixups, squashes, and other
  git-generated messages. You can add your own patterns with the `ignores` config key,
  or disable the defaults with `disable-default-ignores: true`.

- Can I use the old `version` config key?

  Yes. The `version` key is still accepted for backward compatibility, but new config
  files should use `min-version` instead.

## License

All packages are licensed under [MIT License](LICENSE.md)

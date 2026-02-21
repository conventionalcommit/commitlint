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
      - [Precedence](#precedence)
        - [Config](#config-1)
        - [Message](#message)
    - [hook](#hook)
    - [debug](#debug)
  - [Default Config](#default-config)
    - [Commit Types](#commit-types)
  - [Ignore Patterns](#ignore-patterns)
    - [Default Ignore Patterns](#default-ignore-patterns)
    - [Custom Ignore Patterns](#custom-ignore-patterns)
    - [Disabling Default Ignores](#disabling-default-ignores)
  - [Available Rules](#available-rules)
  - [Available Formatters](#available-formatters)
  - [Extensibility](#extensibility)
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

#### Precedence

`commitlint lint` follows below order for `config` and `message`

##### Config

- config file passed to `--config` command-line argument
- `COMMITLINT_CONFIG` env variable
- config file in current directory or git repo root in the below order
  - .commitlint.yml
  - .commitlint.yaml
  - commitlint.yml
  - commitlint.yaml
- [default config](#default-config)

##### Message

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

Commonly used commit types from [Conventional Commit Types](https://github.com/commitizen/conventional-commit-types)

| Type     | Description                                                                      |
| :------- | :------------------------------------------------------------------------------- |
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

| Pattern | Matches |
| :--- | :--- |
| `^Merge pull request #\d+` | GitHub pull request merges |
| `^Merge .+ into .+` | Generic merge (X into Y) |
| `^Merge branch '.+'` | `git merge` branch |
| `^Merge tag '.+'` | `git merge` tag |
| `^Merge remote-tracking branch '.+'` | `git merge` remote-tracking branch |
| `^Merged .+ (in\|into) .+` | Azure DevOps / Bitbucket merged |
| `^Merged PR #?\d+` | Azure DevOps pull request |
| `^(R\|r)evert ` | `git revert` |
| `^(R\|r)eapply ` | `git reapply` |
| `^(amend\|fixup\|squash)! ` | `git commit --fixup/--squash/--amend` |
| `^Automatic merge` | Automatic merges |
| `^Auto-merged .+ into .+` | Auto-merged branches |
| `^Initial commit$` | Initial commit (exact match) |

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

The list of available lint rules

| name                   | argument                 | flags             | description                                   |
| ---------------------- | ------------------------ | ----------------- | --------------------------------------------- |
| header-min-length      | int                      | n/a               | checks the min length of header (first line)  |
| header-max-length      | int                      | n/a               | checks the max length of header (first line)  |
| body-max-line-length   | int                      | n/a               | checks the max length of each line in body    |
| footer-max-line-length | int                      | n/a               | checks the max length of each line in footer  |
| type-enum              | []string                 | n/a               | restrict type to given list of string         |
| scope-enum             | []string                 | allow-empty: bool | restrict scope to given list of string        |
| footer-enum            | []string                 | n/a               | restrict footer token to given list of string |
| type-min-length        | int                      | n/a               | checks the min length of type                 |
| type-max-length        | int                      | n/a               | checks the max length of type                 |
| scope-min-length       | int                      | n/a               | checks the min length of scope                |
| scope-max-length       | int                      | n/a               | checks the max length of scope                |
| description-min-length | int                      | n/a               | checks the min length of description          |
| description-max-length | int                      | n/a               | checks the max length of description          |
| body-min-length        | int                      | n/a               | checks the min length of body                 |
| body-max-length        | int                      | n/a               | checks the max length of body                 |
| footer-min-length      | int                      | n/a               | checks the min length of footer               |
| footer-max-length      | int                      | n/a               | checks the max length of footer               |
| type-charset           | string                   | n/a               | restricts type to given charset               |
| scope-charset          | string                   | n/a               | restricts scope to given charset              |
| footer-type-enum       | []{token, types, values} | n/a               | enforces footer notes for given type          |

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

## Extensibility

## FAQ

- How to have custom config for each repository?

  Place `.commitlint.yaml` file in repo root directory. linter follows [config precedence](#precedence).

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

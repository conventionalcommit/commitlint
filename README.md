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

### Manual

- run `commitlint hook create` to create `.commitlint/hooks` containing git hooks
- To enable in single repo
  - run `git config core.hooksPath /path/to/.commitlint/hooks`
- To enable globally
  - run `git config --global core.hooksPath /path/to/.commitlint/hooks`

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

- To create config file, run `commitlint config create` this will create `commitlint.yaml`

- To validate config file, run `commitlint config check --config=/path/to/conf.yaml`

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

- To create hook files, run `commitlint hook create`

### debug

  To prints useful information for debugging commitlint

  run `commitlint debug`

## Default Config

```yaml
version: v0.8.0
formatter: default
rules:
- header-min-length
- header-max-length
- body-max-line-length
- footer-max-line-length
- type-enum
severity:
  default: error
settings:
  body-max-line-length:
    argument: 72
  footer-max-line-length:
    argument: 72
  header-min-length:
    argument: 10
  header-max-length:
    argument: 50
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

  To create a sample config, run `commitlint config create`

- How can I skip lint check for a commit?

  use `--no-verify` flag with `git commit` which skips commit hooks

## License

All packages are licensed under [MIT License](LICENSE.md)

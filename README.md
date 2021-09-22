# commitlint

commitlint checks if your commit messages meets the [conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/conventionalcommit/commitlint)](https://pkg.go.dev/github.com/conventionalcommit/commitlint)

#### Table of Contents

- [Installation](#installation)
  - [Releases](#releases)
  - [Using go](#using-go)
- [Enable in Git Repo](#enable-in-git-repo)
- [Quick Test](#quick-test)
- [Benefits of using conventional commit](#benefits-of-using-conventional-commit)
  - [Commit Types](#commit-types)
- [Commands](#commands)
    * [Custom config for each repo](#custom-config-for-each-repo)
    * [Verify config file](#verify-config-file)
- [Rules](#rules)
- [Library](#library)
  - [Config Precedence](#config-precedence)
  - [Commit Message Precedence](#commit-message-precedence)
- [Default Config](#default-config)
- [License](#license)

### Installation

#### Releases

Download binary from [releases](https://github.com/conventionalcommit/commitlint/releases) and add it to your `PATH`

#### Using go

```bash
go install github.com/conventionalcommit/commitlint@latest
```

### Enable in Git Repo

- enable for a single repository, `cd` to repository directory

  ```bash
  commitlint init
  ```

- enable globally for all repositories

  ```bash
  commitlint init --global
  ```

### Quick Test

```bash
# invalid commit message
echo "fear: do not fear for commit message" | commitlint lint
# ❌ type-enum: type 'fear' is not allowed, you can use one of [feat fix docs style refactor perf test build ci chore revert merge]

# valid commit message
echo "feat: good commit message" | commitlint lint
# ✔ commit message
```

### Benefits of using conventional commit

- Conventional Commit format. Read [Full Specification here](https://www.conventionalcommits.org/en/v1.0.0/#specification)
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

- [Why Use Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#why-use-conventional-commits)

#### Commit Types

Commonly used commit types from [Conventional Commit Types](https://github.com/commitizen/conventional-commit-types)

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
| merge    | Merges a branch                                                                  |

### Commands

##### Custom config for each repo

- run `commitlint create config` in repo root directory

  this will create `commitlint.yaml` in that directory, you can customise the config to your need

##### Verify config file

- run `commitlint verify` to verify if config is valid or not (according to config precedence)

- run `commitlint verify --config=/path/to/conf.yaml`to verify given config file

### Rules

The list of available lint rules

| name                   | argument | flags             | description                                  |
| ---------------------- | -------- | ----------------- | -------------------------------------------- |
| header-min-length      | int      | n/a               | checks the min length of header (first line) |
| header-max-length      | int      | n/a               | checks the max length of header (first line) |
| body-max-line-length   | int      | n/a               | checks the max length of each line in body   |
| footer-max-line-length | int      | n/a               | checks the max length of each line in footer |
| type-enum              | []string | n/a               | restrict type to given list of string        |
| scope-enum             | []string | allow-empty: bool | restrict scope to given list of string       |
| type-min-length        | int      | n/a               | checks the min length of type                |
| type-max-length        | int      | n/a               | checks the max length of type                |
| scope-min-length       | int      | n/a               | checks the min length of scope               |
| scope-max-length       | int      | n/a               | checks the max length of scope               |
| description-min-length | int      | n/a               | checks the min length of description         |
| description-max-length | int      | n/a               | checks the max length of description         |
| body-min-length        | int      | n/a               | checks the min length of body                |
| body-max-length        | int      | n/a               | checks the max length of body                |
| footer-min-length      | int      | n/a               | checks the min length of footer              |
| footer-max-length      | int      | n/a               | checks the max length of footer              |
| type-charset           | string   | n/a               | restricts type to given charset              |
| scope-charset          | string   | n/a               | restricts scope to given charset             |

### Library

#### Config Precedence

- `commitlint.yaml` config file in current directory
- config file passed with `--config` command-line argument
- [default config](#default-config)

#### Commit Message Precedence

- `stdin` pipe stream
- commit message file passed with `--message` command-line argument
- `.git/COMMIT_EDITMSG` in current directory

### Default Config

```yaml
formatter: default
rules:
    header-min-length:
        enabled: true
        severity: error
        argument: 10
    header-max-length:
        enabled: true
        severity: error
        argument: 50
    body-max-line-length:
        enabled: true
        severity: error
        argument: 72
    footer-max-line-length:
        enabled: true
        severity: error
        argument: 72
    type-enum:
        enabled: true
        severity: error
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
            - merge
```

### License

All packages are licensed under [MIT License](https://github.com/conventionalcommit/commitlint/tree/master/LICENSE.md)

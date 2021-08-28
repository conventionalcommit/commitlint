# commitlint

commitlint checks if your commit messages meets the [conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/conventionalcommit/commitlint)](https://pkg.go.dev/github.com/conventionalcommit/commitlint)

#### Table of Contents

- [Installation](#installation)
  - [Releases](#releases)
  - [Using go](#using-go)
- [Enable in Git Repo](#enable-in-git-repo)
- [Test](#test)
- [Custom config for each repo](#custom-config-for-each-repo)
- [Benefits using commitlint](#benefits-using-commitlint)
- [Commit Types](#commit-types)
- [Library](#library)
  - [Config Precedence](#config-precedence)
  - [Message Precedence](#message-precedence)
  - [Default Config](#default-config)
- [License](#license)

### Installation

#### Releases

Download binary from [releases](https://github.com/conventionalcommit/commitlint/releases) and add the path to your `PATH`

#### Using go

```bash
# Install commitlint
go install github.com/conventionalcommit/commitlint@latest
```

### Enable in Git Repo

```bash
# enable for single repo
commitlint init # from repo directory

# enable for globally for all repos
commitlint init --global
```

### Test

```bash
# invalid commit message
echo "fear: do not fear for commit message" | commitlint lint
# ❌ type-enum: type 'fear' is not allowed, you can use one of [feat fix docs style refactor perf test build ci chore revert merge]

# valid commit message
echo "feat: good commit message" | commitlint lint
# ✔ commit message
```

### Custom config for each repo

- run `commitlint create config` in repo root directory
- this will create `commitlint.yaml` in that directory
- you can customise the config to your need

### Benefits using commitlint

- [Why Use Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#why-use-conventional-commits)

### Commit Types

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

### Library

#### Config Precedence

- `commitlint.yaml` config file in current directory
- config file passed with `--config` command-line argument
- [default config](#default-config)

#### Message Precedence

- `stdin` stream
- commit message file passed with `--message` command-line argument
- `.git/COMMIT_EDITMSG` in current directory

#### Default Config

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

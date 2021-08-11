# commitlint

commitlint checks if your commit messages meets the [conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/)

[![PkgGoDev](https://pkg.go.dev/badge/github.com/conventionalcommit/parser)](https://pkg.go.dev/github.com/conventionalcommit/parser)

### Installation

```bash
# Install commitlint
go install github.com/conventionalcommit/commitlint/cmd/commitlint@latest

# test commitlint - error case
echo "fear: do not fear for commit message" | commitlint lint
# will show error message like
# ❌ type: 'fear' is not allowed, you can use one of [feat fix docs style refactor perf test build ci chore revert merge]

# test commitlint - valid case
echo "feat: good commit message" | commitlint lint
#  ✓ commit message
```

### Enable in Git Repo

```bash
# enable for single repo
commitlint create hook # from repo directory

# enable for globally for all repos
commitlint create hook --global
```

### Benefits using commitlint

- [Why Use Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#why-use-conventional-commits)

### Commit Types

​	Commonly used commit types from [Conventional Commit Types](https://github.com/commitizen/conventional-commit-types)

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
header:
    min-length:
        enabled: true
        type: error
        value: 10
    max-length:
        enabled: true
        type: error
        value: 50
    types:
        enabled: true
        type: error
        value:
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
    scopes:
        enabled: false
        type: error
        value: []
body:
    can-be-empty: true
    max-line-length:
        enabled: true
        type: error
        value: 72
footer:
    can-be-empty: true
    max-line-length:
        enabled: true
        type: error
        value: 72
```

### License

All packages are licensed under [MIT License](https://github.com/conventionalcommit/commitlint/tree/master/LICENSE.md)


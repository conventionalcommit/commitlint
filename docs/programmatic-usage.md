# Programmatic Usage

- [Programmatic Usage](#programmatic-usage)
  - [Key Packages](#key-packages)
  - [One-liner with default config](#one-liner-with-default-config)
  - [Full control with default config](#full-control-with-default-config)
  - [Lint with a config file](#lint-with-a-config-file)
  - [Custom rules](#custom-rules)
  - [Custom formatters](#custom-formatters)
  - [Changelog generation](#changelog-generation)

All public packages - [Programmatic Usage](#programmatic-usage)

```bash
go get github.com/conventionalcommit/commitlint@latest
```

## Key Packages

| Package               | Purpose                                                                |
|:----------------------|:-----------------------------------------------------------------------|
| `config`              | Parse config files, build a `Linter`, access defaults                  |
| `lint`                | Core types: `Linter`, `Rule`, `Formatter`, `Config`, `Result`, `Issue` |
| `registry`            | Register and look up custom rules / formatters                         |
| `lint/rule`           | Built-in rule implementations                                          |
| `lint/formatter`      | Built-in lint formatters (`default`, `json`)                           |
| `changelog`           | Core changelog types, `Config`, `Formatter` interface                  |
| `changelog/formatter` | Changelog formatters (`markdown`, `json`)                              |

## One-liner with default config

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

## Full control with default config

Build the linter yourself for more control (e.g. to swap the formatter):

```go
package main

import (
    "fmt"
    "github.com/conventionalcommit/commitlint/config"
    "github.com/conventionalcommit/commitlint/lint/formatter"
)

func main() {
    conf := config.NewDefault()
    // optionally customise conf.Lint here

    linter, err := config.NewLinter(conf.Lint)
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

## Lint with a config file

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

    linter, err := config.NewLinter(conf.Lint)
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

## Custom rules

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

    conf := config.NewDefaultLint()
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

## Custom formatters

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

    conf := config.NewDefaultLint()
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

## Changelog generation

Generate a changelog programmatically:

```go
package main

import (
    "fmt"
    "github.com/conventionalcommit/commitlint/config"
    "github.com/conventionalcommit/commitlint/registry"
)

func main() {
    // Use default changelog config
    clConf := config.NewDefaultChangelog()

    // Create generator for current directory
    gen, err := config.NewGenerator(".", clConf)
    if err != nil {
        panic(err)
    }

    // Generate changelog for all versions
    cl, err := gen.GenerateAll()
    if err != nil {
        panic(err)
    }

    // Format with markdown
    f, _ := registry.GetChangelogFormatter("markdown")
    result, err := f.Format(cl)
    if err != nil {
        panic(err)
    }
    fmt.Print(result)
}
```

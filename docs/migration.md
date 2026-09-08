# Migration Guide

## Migrating from v0.12.0 (or earlier) to the current version

Starting with the version after v0.12.0, the config file format changed from a **flat structure** to a **nested structure** with `lint:` and `changelog:` top-level keys.

commitlint will detect old-format config files and display a helpful error message pointing you here.

### Config file changes

#### Before (flat format, v0.12.0 and earlier)

```yaml
version: v0.11.0
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
  header-min-length:
    argument: 10
  header-max-length:
    argument: 72
  body-max-line-length:
    argument: 100
  footer-max-line-length:
    argument: 100
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

#### After (nested format, current)

```yaml
lint:
  min-version: v0.12.0
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
    header-min-length:
      argument: 10
    header-max-length:
      argument: 72
    body-max-line-length:
      argument: 100
    footer-max-line-length:
      argument: 100
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
changelog:
  formatter: markdown
  header: "# Changelog"
```

### Step-by-step migration

1. **Wrap lint config under `lint:` key**

   Indent your entire existing config by two spaces and add `lint:` as the first line.

2. **Rename `version` to `min-version`**

   The `version` key is still accepted for backward compatibility, but `min-version` is preferred.

   ```yaml
   # Before
   version: v0.11.0

   # After
   lint:
     min-version: v0.12.0
   ```

3. **Add `changelog:` section (optional)**

   The changelog section is entirely optional. If omitted, sensible defaults are applied:

   ```yaml
   changelog:
     formatter: markdown
     header: "# Changelog"
     issue-prefixes:
       - "#"
     include-breaking: true
     skip-merge-commits: true
     types:
       - type: feat
         header: Features
       - type: fix
         header: Bug Fixes
       # ... see README for full list
   ```

4. **Regenerate config (alternative)**

   Instead of manually editing, you can generate a fresh config file:

   ```bash
   # Compact format (only enabled rules)
   commitlint config create --replace

   # Full format (all rules and settings)
   commitlint config create --replace --all
   ```

   Then re-apply any customizations you had.

### What changed and why

| Aspect              | Old format              | New format                                                |
|:--------------------|:------------------------|:----------------------------------------------------------|
| Top-level structure | Flat — all keys at root | Nested under `lint:` and `changelog:`                     |
| Version key         | `version`               | `min-version` (under `lint:`)                             |
| Changelog config    | Not supported           | `changelog:` section                                      |
| Ignore patterns     | Not supported           | `ignores:` and `disable-default-ignores:` (under `lint:`) |

The nested format allows commitlint to manage both linting and changelog generation from a single config file, with clear separation of concerns.

// Package config contains helpers, defaults for linter and changelog
package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/conventionalcommit/commitlint/changelog"
	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/lint"
)

// Config is the top-level configuration that bundles lint and changelog config.
type Config struct {
	Lint      *lint.Config      `yaml:"lint"`
	Changelog *changelog.Config `yaml:"changelog"`
}

// Parse parses the given config file and returns a Config instance.
func Parse(confPath string) (*Config, error) {
	confPath = filepath.Clean(confPath)
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("config file error: %w", err)
	}

	conf := &Config{
		Lint: &lint.Config{
			Formatter: defaultLintFormatter,
			Severity: lint.SeverityConfig{
				Default: lint.SeverityError,
			},
		},
		Changelog: NewDefaultChangelog(),
	}

	err = yaml.UnmarshalStrict(confBytes, conf)
	if err != nil {
		// Detect old flat config format (pre-v0.12.0) that lacks the top-level "lint:" key.
		if isOldConfigFormat(confBytes) {
			return nil, fmt.Errorf(
				"config file error: this looks like a pre-v0.12.0 config (flat format without 'lint:' key).\n"+
					"Please migrate to the new format. See: https://github.com/conventionalcommit/commitlint/blob/main/docs/migration.md\n"+
					"Original error: %w", err,
			)
		}
		return nil, fmt.Errorf("config file error: %w", err)
	}

	// --- Apply lint defaults ---

	// Backward compatibility: accept old "version" key
	if conf.Lint.MinVersion == "" && conf.Lint.DeprecatedVersion != "" {
		conf.Lint.MinVersion = conf.Lint.DeprecatedVersion
	}
	conf.Lint.DeprecatedVersion = ""

	// Default to current version if neither key was provided
	if conf.Lint.MinVersion == "" {
		conf.Lint.MinVersion = internal.Version()
	}

	// Always set the built-in default patterns
	conf.Lint.DefaultIgnorePatterns = DefaultIgnorePatterns()

	// --- Apply changelog defaults ---
	if conf.Changelog == nil {
		conf.Changelog = NewDefaultChangelog()
	} else {
		defaults := NewDefaultChangelog()
		cl := conf.Changelog

		if cl.Formatter == "" {
			cl.Formatter = defaults.Formatter
		}
		if cl.Header == "" {
			cl.Header = defaults.Header
		}
		if len(cl.IssuePrefixes) == 0 {
			cl.IssuePrefixes = defaults.IssuePrefixes
		}
		if len(cl.Types) == 0 {
			cl.Types = defaults.Types
		}
	}

	// --- Validate essentials ---

	if conf.Lint.Formatter == "" {
		return nil, errors.New("config error: lint formatter is empty")
	}

	err = isValidVersion(conf.Lint.MinVersion)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// isOldConfigFormat checks if the YAML bytes look like the old flat config
// (pre-v0.12.0) that had top-level keys like "formatter:", "rules:", "settings:"
// instead of being nested under "lint:".
func isOldConfigFormat(data []byte) bool {
	// Quick heuristic: try to unmarshal into a map and check for old top-level keys
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return false
	}
	// Old format had these at the top level
	oldKeys := []string{"formatter", "rules", "settings", "severity"}
	matches := 0
	for _, k := range oldKeys {
		if _, ok := raw[k]; ok {
			matches++
		}
	}
	// If the file has "lint:" key, it's the new format (even if malformed)
	_, hasLint := raw["lint"]
	return !hasLint && matches >= 2
}

// LookupAndParse gets the config path according to the precedence,
// parses the config file if found, and returns a Config instance.
func LookupAndParse() (*Config, error) {
	confFilePath, confType, err := internal.LookupConfigPath()
	if err != nil {
		return nil, err
	}

	if confType == internal.DefaultConfig {
		return NewDefault(), nil
	}

	conf, err := Parse(confFilePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// WriteTo writes config in yaml format to given io.Writer, including all
// settings and every field even if empty or zero-valued.
func WriteTo(w io.Writer, conf *Config) (retErr error) {
	out := prepareForWrite(conf)
	enc := yaml.NewEncoder(w)
	defer func() {
		err := enc.Close()
		if retErr == nil && err != nil {
			retErr = err
		}
	}()
	return enc.Encode(out)
}

// WriteCompactTo writes config in yaml format to given io.Writer.
// Only settings for enabled rules and non-hidden changelog types are written,
// keeping the output compact.
func WriteCompactTo(w io.Writer, conf *Config) error {
	out := prepareForWrite(conf)

	// Only settings for enabled lint rules
	if out.Lint != nil && len(out.Lint.Rules) > 0 && len(out.Lint.Settings) > 0 {
		enabled := make(map[string]struct{}, len(out.Lint.Rules))
		for _, r := range out.Lint.Rules {
			enabled[r] = struct{}{}
		}
		filtered := make(map[string]lint.RuleSetting, len(out.Lint.Rules))
		for name, setting := range out.Lint.Settings {
			if _, ok := enabled[name]; ok {
				filtered[name] = setting
			}
		}
		out.Lint.Settings = filtered
	}

	// Only non-hidden changelog types
	if out.Changelog != nil && len(out.Changelog.Types) > 0 {
		visible := make([]changelog.TypeConfig, 0, len(out.Changelog.Types))
		for _, tc := range out.Changelog.Types {
			if !tc.Hidden {
				visible = append(visible, tc)
			}
		}
		// Shallow-copy to avoid mutating the caller's config
		clCopy := *out.Changelog
		clCopy.Types = visible
		out.Changelog = &clCopy
	}

	enc := yaml.NewEncoder(w)
	defer enc.Close()
	return enc.Encode(out)
}

// prepareForWrite returns a copy with nil pointers filled with defaults.
func prepareForWrite(conf *Config) *Config {
	out := *conf
	if out.Lint == nil {
		out.Lint = NewDefaultLint()
	}
	if out.Changelog == nil {
		out.Changelog = NewDefaultChangelog()
	}
	return &out
}

func Validate(conf *Config) []error {
	var errs []error
	if conf.Lint != nil {
		errs = append(errs, ValidateLint(conf.Lint)...)
	} else {
		errs = append(errs, errors.New("lint config is nil"))
	}
	if conf.Changelog != nil {
		errs = append(errs, ValidateChangelog(conf.Changelog)...)
	} else {
		errs = append(errs, errors.New("changelog config is nil"))
	}
	return errs
}

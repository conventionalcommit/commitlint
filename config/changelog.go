package config

import (
	"fmt"

	"github.com/conventionalcommit/commitlint/changelog"
	"github.com/conventionalcommit/commitlint/registry"
)

// NewGenerator creates a changelog generator for the given repository directory and config.
func NewGenerator(repoDir string, conf *changelog.Config) (*changelog.Generator, error) {
	return changelog.New(repoDir, conf)
}

// GetChangelogFormatter returns the changelog formatter as defined in conf.
func GetChangelogFormatter(conf *changelog.Config) (changelog.Formatter, error) {
	if conf.Formatter == "" {
		return nil, fmt.Errorf("config error: changelog formatter is empty")
	}
	f, ok := registry.GetChangelogFormatter(conf.Formatter)
	if !ok {
		return nil, fmt.Errorf("config error: '%s' changelog formatter not found", conf.Formatter)
	}
	return f, nil
}

// ValidateChangelog validates the given changelog config.
// It checks that the formatter is registered and types are defined.
func ValidateChangelog(conf *changelog.Config) []error {
	var errs []error

	if conf.Formatter == "" {
		errs = append(errs, fmt.Errorf("changelog formatter is empty"))
	} else {
		_, ok := registry.GetChangelogFormatter(conf.Formatter)
		if !ok {
			errs = append(errs, fmt.Errorf("unknown changelog formatter '%s'", conf.Formatter))
		}
	}

	if len(conf.Types) == 0 {
		errs = append(errs, fmt.Errorf("changelog types are empty"))
	}

	// Check for duplicate types
	seen := make(map[string]struct{}, len(conf.Types))
	for _, tc := range conf.Types {
		if tc.Type == "" {
			errs = append(errs, fmt.Errorf("changelog type entry has empty type field"))
			continue
		}
		if tc.Header == "" {
			errs = append(errs, fmt.Errorf("changelog type '%s' has empty header", tc.Type))
		}
		if _, exists := seen[tc.Type]; exists {
			errs = append(errs, fmt.Errorf("duplicate changelog type '%s'", tc.Type))
		} else {
			seen[tc.Type] = struct{}{}
		}
	}

	return errs
}

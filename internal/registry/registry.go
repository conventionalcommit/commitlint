// Package registry contains registered rules and formatters
package registry

import (
	"fmt"
	"sync"

	"github.com/conventionalcommit/commitlint/formatter"
	"github.com/conventionalcommit/commitlint/lint"
	"github.com/conventionalcommit/commitlint/rule"
)

var globalRegistry = newRegistry()

// RegisterRule registers a custom rule
// if rule already exists, returns error
func RegisterRule(r lint.Rule) error {
	return globalRegistry.RegisterRule(r)
}

// RegisterFormatter registers a custom formatter
// if formatter already exists, returns error
func RegisterFormatter(format lint.Formatter) error {
	return globalRegistry.RegisterFormatter(format)
}

// GetRule returns Rule with given name
func GetRule(name string) (lint.Rule, bool) {
	return globalRegistry.GetRule(name)
}

// GetFormatter returns Formatter with given name
func GetFormatter(name string) (lint.Formatter, bool) {
	return globalRegistry.GetFormatter(name)
}

// Rules returns all registered rules
func Rules() []lint.Rule {
	return globalRegistry.Rules()
}

type registry struct {
	mut *sync.Mutex

	allRules      map[string]lint.Rule
	allFormatters map[string]lint.Formatter
}

func newRegistry() *registry {
	defaultRules := []lint.Rule{
		// Length rules
		&rule.BodyMinLenRule{}, &rule.BodyMaxLenRule{},
		&rule.FooterMinLenRule{}, &rule.FooterMaxLenRule{},
		&rule.HeadMaxLenRule{}, &rule.HeadMinLenRule{},
		&rule.BodyMaxLineLenRule{}, &rule.FooterMaxLineLenRule{},
		&rule.TypeMaxLenRule{}, &rule.ScopeMaxLenRule{}, &rule.DescriptionMaxLenRule{},
		&rule.TypeMinLenRule{}, &rule.ScopeMinLenRule{}, &rule.DescriptionMinLenRule{},

		// Enum rules
		&rule.TypeEnumRule{}, &rule.ScopeEnumRule{}, &rule.FooterEnumRule{},
		&rule.FooterTypeEnumRule{},

		// Charset rules
		&rule.TypeCharsetRule{}, &rule.ScopeCharsetRule{},

		// Case rules
		&rule.TypeCaseRule{}, &rule.ScopeCaseRule{},
		&rule.DescriptionCaseRule{}, &rule.BodyCaseRule{}, &rule.HeaderCaseRule{},

		// Empty rules
		&rule.TypeEmptyRule{}, &rule.ScopeEmptyRule{},
		&rule.BodyEmptyRule{}, &rule.FooterEmptyRule{}, &rule.DescriptionEmptyRule{},

		// Full-stop rules
		&rule.HeaderFullStopRule{}, &rule.BodyFullStopRule{}, &rule.DescriptionFullStopRule{},

		// Leading-blank rules
		&rule.BodyLeadingBlankRule{}, &rule.FooterLeadingBlankRule{},

		// Header trim
		&rule.HeaderTrimRule{},

		// Trailer / signed-off-by rules
		&rule.SignedOffByRule{}, &rule.TrailerExistsRule{},

		// Breaking change
		&rule.BreakingChangeExclamationMarkRule{},
	}

	defaultFormatters := []lint.Formatter{
		&formatter.DefaultFormatter{},
		&formatter.JSONFormatter{},
	}

	reg := &registry{
		mut: &sync.Mutex{},

		allRules:      make(map[string]lint.Rule),
		allFormatters: make(map[string]lint.Formatter),
	}

	// Register Default Rules
	for _, cRule := range defaultRules {
		err := reg.RegisterRule(cRule)
		if err != nil {
			// default rules should not throw error
			panic(err)
		}
	}

	// Register Default Formatters
	for _, format := range defaultFormatters {
		err := reg.RegisterFormatter(format)
		if err != nil {
			// default formatters should not throw error
			panic(err)
		}
	}

	return reg
}

func (reg *registry) RegisterRule(cRule lint.Rule) error {
	reg.mut.Lock()
	defer reg.mut.Unlock()

	_, ok := reg.allRules[cRule.Name()]
	if ok {
		return fmt.Errorf("'%s' rule already registered", cRule.Name())
	}

	reg.allRules[cRule.Name()] = cRule

	return nil
}

func (reg *registry) RegisterFormatter(format lint.Formatter) error {
	reg.mut.Lock()
	defer reg.mut.Unlock()

	_, ok := reg.allFormatters[format.Name()]
	if ok {
		return fmt.Errorf("'%s' formatter already registered", format.Name())
	}

	reg.allFormatters[format.Name()] = format

	return nil
}

func (reg *registry) GetRule(name string) (lint.Rule, bool) {
	reg.mut.Lock()
	defer reg.mut.Unlock()

	cRule, ok := reg.allRules[name]
	return cRule, ok
}

func (reg *registry) GetFormatter(name string) (lint.Formatter, bool) {
	reg.mut.Lock()
	defer reg.mut.Unlock()

	format, ok := reg.allFormatters[name]
	return format, ok
}

func (reg *registry) Formatters() []lint.Formatter {
	reg.mut.Lock()
	defer reg.mut.Unlock()

	allFormats := make([]lint.Formatter, 0, len(reg.allFormatters))
	for _, f := range reg.allFormatters {
		allFormats = append(allFormats, f)
	}
	return allFormats
}

func (reg *registry) Rules() []lint.Rule {
	reg.mut.Lock()
	defer reg.mut.Unlock()

	allRules := make([]lint.Rule, 0, len(reg.allRules))
	for _, r := range reg.allRules {
		allRules = append(allRules, r)
	}
	return allRules
}

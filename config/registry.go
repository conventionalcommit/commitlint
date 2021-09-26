package config

import (
	"fmt"
	"sync"

	"github.com/conventionalcommit/commitlint/lint"
)

var globalRegistry = newRegistry()

// RegisterRule registers a custom rule
// if rule already exists, returns error
func RegisterRule(rule lint.Rule) error {
	return globalRegistry.RegisterRule(rule)
}

// RegisterFormatter registers a custom formatter
// if formatter already exists, returns error
func RegisterFormatter(format lint.Formatter) error {
	return globalRegistry.RegisterFormatter(format)
}

type registry struct {
	mut *sync.Mutex

	allRules      map[string]lint.Rule
	allFormatters map[string]lint.Formatter
}

func newRegistry() *registry {
	reg := &registry{
		mut: &sync.Mutex{},

		allRules:      make(map[string]lint.Rule),
		allFormatters: make(map[string]lint.Formatter),
	}

	// Register Default Rules
	for _, rule := range defaultRules {
		err := reg.RegisterRule(rule)
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

func (r *registry) RegisterRule(rule lint.Rule) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	_, ok := r.allRules[rule.Name()]
	if ok {
		return fmt.Errorf("'%s' rule already registered", rule.Name())
	}

	r.allRules[rule.Name()] = rule

	return nil
}

func (r *registry) RegisterFormatter(format lint.Formatter) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	_, ok := r.allFormatters[format.Name()]
	if ok {
		return fmt.Errorf("'%s' formatter already registered", format.Name())
	}

	r.allFormatters[format.Name()] = format

	return nil
}

func (r *registry) GetRule(name string) (lint.Rule, bool) {
	r.mut.Lock()
	defer r.mut.Unlock()

	rule, ok := r.allRules[name]
	return rule, ok
}

func (r *registry) GetFormatter(name string) (lint.Formatter, bool) {
	r.mut.Lock()
	defer r.mut.Unlock()

	format, ok := r.allFormatters[name]
	return format, ok
}

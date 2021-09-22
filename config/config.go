// Package config contains helpers, defaults for linter
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/conventionalcommit/commitlint/lint"
)

const (
	// ConfFileName represent config file name
	ConfFileName = "commitlint.yaml"
)

// GetConfig returns parses config file and returns Config instance
func GetConfig(flagConfPath string) (*lint.Config, error) {
	confFilePath, useDefault, err := GetConfigPath(flagConfPath)
	if err != nil {
		return nil, err
	}

	if useDefault {
		return defConf, nil
	}

	conf, err := Parse(confFilePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// GetConfigPath returns config file path, follwing below
// 	1. check for conf in current directory
// 	2. check for conf flag
// 	3. load default conf
func GetConfigPath(confFilePath string) (string, bool, error) {
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", false, err
	}

	// check if conf file exists in current directory
	currentDirConf := filepath.Join(currentDir, ConfFileName)
	if _, err1 := os.Stat(currentDirConf); !os.IsNotExist(err1) {
		return currentDirConf, false, nil
	}

	// if confFilePath empty,
	// means no config in current directory or config flag is empty
	// use default config
	if confFilePath == "" {
		return "", true, nil
	}
	return filepath.Clean(confFilePath), false, nil
}

// Parse parse Config from given file
func Parse(confPath string) (*lint.Config, error) {
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	conf := &lint.Config{}
	err = yaml.Unmarshal(confBytes, conf)
	if err != nil {
		return nil, err
	}

	err = Validate(conf)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return conf, nil
}

// Validate parses Config from given data
func Validate(conf *lint.Config) error {
	if conf.Formatter == "" {
		return errors.New("formatter is empty")
	}

	_, ok := globalRegistry.GetFormatter(conf.Formatter)
	if !ok {
		return fmt.Errorf("unknown formatter '%s'", conf.Formatter)
	}

	for ruleName, r := range conf.Rules {
		// Check Severity Level of rule config
		switch r.Severity {
		case lint.SeverityError:
		case lint.SeverityWarn:
		default:
			return fmt.Errorf("unknown severity level '%s' for rule '%s'", r.Severity, ruleName)
		}

		// Check if rule is registered
		_, ok := globalRegistry.GetRule(ruleName)
		if !ok {
			return fmt.Errorf("unknown rule '%s'", ruleName)
		}
	}

	return nil
}

// DefaultConfToFile writes default config to given file
func DefaultConfToFile(isOnlyEnabled bool) error {
	outPath := filepath.Join(".", filepath.Clean(ConfFileName))
	if isOnlyEnabled {
		confClone := &lint.Config{
			Formatter: defConf.Formatter,
			Rules:     map[string]lint.RuleConfig{},
		}

		for ruleName, r := range defConf.Rules {
			if r.Enabled {
				confClone.Rules[ruleName] = r
			}
		}

		return WriteConfToFile(outPath, confClone)
	}
	return WriteConfToFile(outPath, defConf)
}

// WriteConfToFile util func to write config object to given file
func WriteConfToFile(outFilePath string, conf *lint.Config) (retErr error) {
	file, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	w := bufio.NewWriter(file)
	defer func() {
		err1 := w.Flush()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	enc := yaml.NewEncoder(w)
	return enc.Encode(conf)
}

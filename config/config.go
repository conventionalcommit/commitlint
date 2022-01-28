// Package config contains helpers, defaults for linter
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/semver"
	yaml "gopkg.in/yaml.v3"

	"github.com/conventionalcommit/commitlint/internal"
	"github.com/conventionalcommit/commitlint/internal/registry"
	"github.com/conventionalcommit/commitlint/lint"
)

const (
	// ConfigFile represent default config file name
	ConfigFile = "commitlint.yaml"
)

// GetConfig gets the config path according to the precedence
// if needed parses given config file and returns config instance
func GetConfig(confPath string) (*lint.Config, error) {
	confFilePath, useDefault, err := getConfigPath(confPath)
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

	if conf.Formatter == "" {
		return nil, errors.New("conf error: formatter is empty")
	}

	err = isValidVersion(conf.MinVersion)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// getConfigPath returns config file path following below order
// 	1. commitlint.yaml in current directory
// 	2. confFilePath parameter
// 	3. use default config
func getConfigPath(confFilePath string) (confPath string, isDefault bool, retErr error) {
	// get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", false, err
	}

	// check if conf file exists in current directory
	currentDirConf := filepath.Join(currentDir, ConfigFile)
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

// Parse parse given file in confPath, and return Config instance, error if any
func Parse(confPath string) (*lint.Config, error) {
	confPath = filepath.Clean(confPath)
	confBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	conf := &lint.Config{}
	err = yaml.Unmarshal(confBytes, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Validate validates given config instance, it checks the following
// If formatters, rules are registered/known
// If arguments to rules are valid
// If version is valid and atleast minimum than commitlint version used
func Validate(conf *lint.Config) []error {
	var errs []error

	if conf.Formatter == "" {
		errs = append(errs, errors.New("formatter is empty"))
	} else {
		_, ok := registry.GetFormatter(conf.Formatter)
		if !ok {
			errs = append(errs, fmt.Errorf("unknown formatter '%s'", conf.Formatter))
		}
	}

	err := isValidVersion(conf.MinVersion)
	if err != nil {
		errs = append(errs, err)
	}

	for ruleName, r := range conf.Rules {
		// Check Severity Level of rule config
		switch r.Severity {
		case lint.SeverityError:
		case lint.SeverityWarn:
		default:
			errs = append(errs, fmt.Errorf("unknown severity level '%s' for rule '%s'", r.Severity, ruleName))
		}

		// Check if rule is registered
		ruleData, ok := registry.GetRule(ruleName)
		if !ok {
			errs = append(errs, fmt.Errorf("unknown rule '%s'", ruleName))
			continue
		}

		err := ruleData.Apply(r.Argument, r.Flags)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// WriteToFile util func to write config object to given file
func WriteToFile(outFilePath string, conf *lint.Config) (retErr error) {
	f, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if retErr == nil && err != nil {
			retErr = err
		}
	}()

	w := bufio.NewWriter(f)
	defer func() {
		err := w.Flush()
		if retErr == nil && err != nil {
			retErr = err
		}
	}()

	enc := yaml.NewEncoder(w)
	return enc.Encode(conf)
}

func isValidVersion(versionNo string) error {
	if versionNo == "" {
		return errors.New("version is empty")
	}
	if !semver.IsValid(versionNo) {
		return errors.New("invalid version should be in semver format")
	}
	return nil
}

func checkIfMinVersion(versionNo string) error {
	cmp := semver.Compare(internal.Version(), versionNo)
	if cmp != -1 {
		return nil
	}
	return fmt.Errorf("min version required is %s. you have %s.\nupgrade commitlint", versionNo, internal.Version())
}

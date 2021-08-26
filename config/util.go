package config

import (
	"bufio"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/conventionalcommit/commitlint/lint"
)

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

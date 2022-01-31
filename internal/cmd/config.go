package cmd

import (
	"os"
	"path/filepath"

	"github.com/conventionalcommit/commitlint/config"
)

// configCreate is the callback function for create config command
func configCreate(fileName string, isReplace bool) error {
	defConf := config.Default()
	outPath := filepath.Join(".", fileName)
	// if config file already exists skip creating or overwriting it
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		if !isReplace {
			return errConfigExist
		}
	}
	return config.WriteToFile(outPath, defConf)
}

// configCheck is the callback function for check/verify command
func configCheck(confPath string) []error {
	conf, err := config.Parse(confPath)
	if err != nil {
		return []error{err}
	}
	return config.Validate(conf)
}

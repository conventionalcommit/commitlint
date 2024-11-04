package cmd

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/conventionalcommit/commitlint/config"
)

// configCreate is the callback function for create config command
func configCreate(fileName string, isReplace bool) (retErr error) {
	outPath := filepath.Join(".", fileName)

	// if config file already exists skip creating or overwriting it
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		if !isReplace {
			return handleError(errConfigExist, "Config file already exists")
		}
	}

	outFilePath := filepath.Clean(outPath)
	f, err := os.Create(outFilePath)
	if handleError(err, "Failed to create config file") != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if retErr == nil && err != nil {
			retErr = handleError(err, "Failed to close config file")
		}
	}()

	w := bufio.NewWriter(f)
	defer func() {
		err := w.Flush()
		if retErr == nil && err != nil {
			retErr = handleError(err, "Failed to flush writer")
		}
	}()

	defConf := config.NewDefault()
	return handleError(config.WriteTo(w, defConf), "Failed to write config to file")
}

// configCheck is the callback function for check/verify command
func configCheck(confPath string) []error {
	conf, err := config.Parse(confPath)
	if handleError(err, "Failed to parse configuration file") != nil {
		return []error{err}
	}
	return config.Validate(conf)
}

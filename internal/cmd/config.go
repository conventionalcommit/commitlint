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
			return errConfigExist
		}
	}

	outFilePath := filepath.Clean(outPath)
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

	defConf := config.NewDefault()
	return config.WriteTo(w, defConf)
}

// configCheck is the callback function for check/verify command
func configCheck(confPath string) []error {
	conf, err := config.Parse(confPath)
	if err != nil {
		return []error{err}
	}
	return config.Validate(conf)
}

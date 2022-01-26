package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	cli "github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint/config"
)

const (
	// errExitCode represent error exit code
	errExitCode = 1
)

// initLint is the callback function for init command
func initLint(confPath string, isGlobal, isReplace bool) error {
	hookDir, err := initHooks(confPath, isGlobal, isReplace)
	if err != nil {
		return err
	}
	return setGitConf(hookDir, isGlobal)
}

// lintMsg is the callback function for lint command
func lintMsg(confPath, msgPath string) error {
	// NOTE: lint should return with exit code for error case
	resStr, hasError, err := runLint(confPath, msgPath)
	if err != nil {
		return cli.Exit(err, errExitCode)
	}

	if hasError {
		return cli.Exit(resStr, errExitCode)
	}

	// print success message
	fmt.Println(resStr)
	return nil
}

// hookCreate is the callback function for create hook command
func hookCreate(confPath string, isReplace bool) error {
	return createHooks(confPath, isReplace)
}

// configCreate is the callback function for create config command
func configCreate(onlyEnabled, isReplace bool) error {
	defConf := config.GetDefaultConfig(onlyEnabled)
	outPath := filepath.Join(".", config.ConfigFile)
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

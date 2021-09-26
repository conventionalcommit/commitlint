package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint/config"
)

const (
	// ErrExitCode represent error exit code
	ErrExitCode = 1
)

// Init is the callback function for init command
func Init(confPath string, isGlobal, isReplace bool) error {
	hookDir, err := initHooks(confPath, isGlobal, isReplace)
	if err != nil {
		return err
	}
	return setGitConf(hookDir, isGlobal)
}

// Lint is the callback function for lint command
func Lint(confPath, msgPath string) error {
	// NOTE: lint should return with exit code for error case
	resStr, hasError, err := runLint(confPath, msgPath)
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}

	if hasError {
		return cli.Exit(resStr, ErrExitCode)
	}

	// print success message
	fmt.Println(resStr)
	return nil
}

// HookCreate is the callback function for create hook command
func HookCreate(confPath string, isReplace bool) error {
	return createHooks(confPath, isReplace)
}

// ConfigCreate is the callback function for create config command
func ConfigCreate(onlyEnabled bool) error {
	return config.DefaultConfToFile(onlyEnabled)
}

// ConfigCheck is the callback function for check/verify command
func ConfigCheck(confPath string) error {
	conf, err := config.Parse(confPath)
	if err != nil {
		return err
	}
	return config.Validate(conf)
}

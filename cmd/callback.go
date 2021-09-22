package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/hook"
)

const (
	// ErrExitCode represent error exit code
	ErrExitCode = 1

	// HookDir represent default hook directory
	HookDir = ".commitlint/hooks"
)

// Init is the callback function for init command
func Init(isGlobal bool) error {
	hookDir, err := getHookDir(isGlobal)
	if err != nil {
		return err
	}

	err = os.MkdirAll(hookDir, os.ModePerm)
	if err != nil {
		return err
	}

	// create hook file
	err = hook.WriteToFile(hookDir)
	if err != nil {
		return err
	}

	err = setGitConf(hookDir, isGlobal)
	if err != nil {
		return err
	}

	fmt.Println("commitlint init successfully")
	return nil
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

// CreateHook is the callback function for create hook command
func CreateHook() (retErr error) {
	return hook.WriteToFile(".")
}

// CreateConfig is the callback function for create config command
func CreateConfig(onlyEnabled bool) error {
	return config.DefaultConfToFile(onlyEnabled)
}

// VerifyConfig is the callback function for verify command
func VerifyConfig(confFlag string) error {
	confPath, useDefault, err := config.GetConfigPath(confFlag)
	if err != nil {
		return err
	}

	if useDefault {
		fmt.Println("no config file found, default config will be used")
		return nil
	}

	_, _, err = getLinter(confPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s config is valid\n", confPath)
	return nil
}

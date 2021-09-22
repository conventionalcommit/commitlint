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

func initCallback(ctx *cli.Context) (retErr error) {
	isGlobal := ctx.Bool("global")

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

func lintCallback(ctx *cli.Context) error {
	confFilePath := ctx.String("config")
	fileInput := ctx.String("message")

	resStr, hasError, err := runLint(confFilePath, fileInput)
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}

	if hasError {
		return cli.Exit(resStr, ErrExitCode)
	}

	fmt.Println(resStr)
	return nil
}

func hookCreateCallback(ctx *cli.Context) (retErr error) {
	err := hook.WriteToFile(".")
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}
	return nil
}

func configCreateCallback(ctx *cli.Context) error {
	isOnlyEnabled := ctx.Bool("enabled")
	err := config.DefaultConfToFile(isOnlyEnabled)
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}
	return nil
}

func verifyCallback(ctx *cli.Context) error {
	confFlag := ctx.String("config")

	confPath, useDefault, err := config.GetConfigPath(confFlag)
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}

	if useDefault {
		fmt.Println("no config file found, default config will be used")
		return nil
	}

	_, _, err = getLinter(confPath)
	if err != nil {
		return cli.Exit(err, ErrExitCode)
	}

	fmt.Printf("%s config is valid\n", confPath)
	return nil
}

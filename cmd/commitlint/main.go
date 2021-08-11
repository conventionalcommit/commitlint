package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/conventionalcommit/commitlint"
)

const (
	exitCode = 1

	defHookDir      = ".commitlint/hooks"
	defConfFileName = "commitlint.yaml"
	commitMsgHook   = "commit-msg"
)

func main() {
	app := getApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func initConfCallback(ctx *cli.Context) error {
	err := commitlint.DefaultConfToFile(defConfFileName)
	if err != nil {
		return cli.Exit(err, exitCode)
	}
	return nil
}

func initHookCallback(ctx *cli.Context) (retErr error) {
	// get user home dir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// create hooks dir
	hookDir := filepath.Join(homeDir, defHookDir)
	err = os.MkdirAll(hookDir, os.ModePerm)
	if err != nil {
		return err
	}

	// create hook file
	hookFile := filepath.Join(hookDir, commitMsgHook)
	err = writeHookFile(hookFile)
	if err != nil {
		return err
	}

	isGlobal := ctx.Bool("global")
	return setGitConf(hookDir, isGlobal)
}

func lintCallback(ctx *cli.Context) error {
	confFilePath := ctx.String("config")
	linter, err := getLinter(confFilePath)
	if err != nil {
		return cli.Exit(err, exitCode)
	}

	commitMsg := readStdIn()
	if commitMsg == "" {
		fileInput := ctx.String("message")
		msg, err2 := readCommitMsg(fileInput)
		if err2 != nil {
			return cli.Exit(err2, exitCode)
		}
		commitMsg = msg
	}

	lintMsg, hasError, err := linter.Lint(commitMsg)
	if err != nil {
		return cli.Exit(err, exitCode)
	}

	if hasError {
		return cli.Exit(lintMsg, exitCode)
	}
	fmt.Println(lintMsg)
	return nil
}

func setGitConf(hookDir string, isGlobal bool) error {
	var args = []string{"config"}
	if isGlobal {
		args = append(args, "--global")
	}
	args = append(args, "core.hooksPath", hookDir)

	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func writeHookFile(hookFilePath string) (retErr error) {
	// if commit-msg already exists skip creating or overwriting it
	if isFileExists(hookFilePath) {
		return nil
	}
	// commit-msg needs to be executable
	file, err := os.OpenFile(hookFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	commitHook := commitlint.CommitMsgHook()
	_, err = file.WriteString(commitHook)
	return err
}

func getLinter(confFilePath string) (*commitlint.Linter, error) {
	// Config Precedence
	// 	1. Check for conf in current directory
	// 	2. Check for conf flag
	// 	3. Load default conf

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// check if conf file exists in current directory
	currentDirConf := filepath.Join(currentDir, defConfFileName)
	if isFileExists(currentDirConf) {
		confFilePath = currentDirConf
	}

	if confFilePath != "" {
		conf, err := readConf(confFilePath)
		if err != nil {
			return nil, err
		}
		return commitlint.NewLinter(conf), nil
	}
	return commitlint.NewDefaultLinter(), nil
}

func readConf(confPath string) (*commitlint.Config, error) {
	confFile, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	conf := &commitlint.Config{}
	err = yaml.Unmarshal(confFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func readCommitMsg(fileInput string) (string, error) {
	// TODO: check if currentDir is inside git repo?
	if fileInput == "" {
		fileInput = "./.git/COMMIT_EDITMSG"
	}

	inBytes, err := os.ReadFile(fileInput)
	if err != nil {
		return "", err
	}
	return string(inBytes), nil
}

func readStdIn() string {
	readBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		// TODO: handle error?
		return ""
	}
	s := string(readBytes)
	return strings.TrimSpace(s)
}

func isFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

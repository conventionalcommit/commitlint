package cmd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/conventionalcommit/commitlint/internal/hook"
)

const (
	defaultHooksPath = ".commitlint/hooks"
)

var (
	errHooksExist  = errors.New("hooks already exists")
	errConfigExist = errors.New("config file already exists")
)

// hookCreate is the callback function for create hook command
func hookCreate(hooksPath string, isReplace bool) error {
	if hooksPath == "" {
		hooksPath = filepath.Join(".", defaultHooksPath)
	}
	hooksPath = filepath.Clean(hooksPath)
	return createHooks(hooksPath, isReplace)
}

func initHooks(confPath, hookFlag string, isGlobal, isReplace bool) (string, error) {
	hookDir, err := getHookDir(hookFlag, isGlobal)
	if err != nil {
		return "", err
	}

	err = writeHooks(hookDir, isReplace)
	if err != nil {
		return "", err
	}
	return hookDir, nil
}

func createHooks(hookBaseDir string, isReplace bool) error {
	return writeHooks(hookBaseDir, isReplace)
}

func writeHooks(hookDir string, isReplace bool) error {
	// if commit-msg already exists skip creating or overwriting it
	if _, err := os.Stat(hookDir); !os.IsNotExist(err) {
		if !isReplace {
			return errHooksExist
		}
	}

	err := os.MkdirAll(hookDir, os.ModePerm)
	if err != nil {
		return err
	}

	// create hook file
	return hook.WriteHooks(hookDir)
}

func getHookDir(hookFlag string, isGlobal bool) (string, error) {
	if hookFlag != "" {
		return filepath.Abs(hookFlag)
	}

	hookFlag = defaultHooksPath

	if isGlobal {
		// get user home dir
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		// create hooks dir
		hookDir := filepath.Join(homeDir, hookFlag)
		return hookDir, nil
	}

	gitDir, err := getRepoRootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir, hookFlag), nil
}

func getRepoRootDir() (string, error) {
	byteOut := &bytes.Buffer{}

	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stdout = byteOut
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	gitDir := filepath.Clean(byteOut.String())

	// remove /.git at last
	gitDir = filepath.Dir(gitDir)

	return gitDir, nil
}

func isHookExists(err error) bool {
	return err == errHooksExist
}

func isConfExists(err error) bool {
	return err == errConfigExist
}

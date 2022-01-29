package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/conventionalcommit/commitlint/internal/hook"
)

const (
	// hookBaseDir represent default hook directory
	hookBaseDir = ".commitlint/hooks"
)

var errHooksExist = errors.New("hooks already exists")
var errConfigExist = errors.New("config file already exists")

func initHooks(confPath string, isGlobal, isReplace bool) (string, error) {
	hookDir, err := getHookDir(hookBaseDir, isGlobal)
	if err != nil {
		return "", err
	}

	err = writeHooks(hookDir, isReplace)
	if err != nil {
		return "", err
	}
	return hookDir, nil
}

func createHooks(isReplace bool) error {
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

func getHookDir(baseDir string, isGlobal bool) (string, error) {
	baseDir = filepath.Clean(baseDir)

	if isGlobal {
		// get user home dir
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		// create hooks dir
		hookDir := filepath.Join(homeDir, baseDir)
		return hookDir, nil
	}

	gitDir, err := getRepoRootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir, baseDir), nil
}

func isHookExists(err error) bool {
	return err == errHooksExist
}

func isConfExists(err error) bool {
	return err == errConfigExist
}

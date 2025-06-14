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
	return handleError(createHooks(hooksPath, isReplace), "Failed to create hooks")
}

func initHooks(confPath, hookFlag string, isGlobal, isReplace bool) (string, error) {
	hookDir, err := getHookDir(hookFlag, isGlobal)
	if handleError(err, "Failed to get hook directory") != nil {
		return "", err
	}

	err = writeHooks(hookDir, isReplace)
	if handleError(err, "Failed to write hooks") != nil {
		return "", err
	}
	return hookDir, nil
}

func createHooks(hookBaseDir string, isReplace bool) error {
	return handleError(writeHooks(hookBaseDir, isReplace), "Failed to write hooks to base directory")
}

func writeHooks(hookDir string, isReplace bool) error {
	// if commit-msg already exists skip creating or overwriting it
	if _, err := os.Stat(hookDir); !os.IsNotExist(err) {
		if !isReplace {
			return handleError(errHooksExist, "Hook already exists and replace option not set")
		}
	}

	err := os.MkdirAll(hookDir, os.ModePerm)
	if handleError(err, "Failed to create hook directory") != nil {
		return err
	}

	// create hook file
	return handleError(hook.WriteHooks(hookDir), "Failed to write hooks to directory")
}

func getHookDir(hookFlag string, isGlobal bool) (string, error) {
	if hookFlag != "" {
		absPath, err := filepath.Abs(hookFlag)
		if handleError(err, "Failed to get absolute path for hook directory") != nil {
			return "", err
		}
		return absPath, nil
	}

	hookFlag = defaultHooksPath

	if isGlobal {
		// get user home dir
		homeDir, err := os.UserHomeDir()
		if handleError(err, "Failed to get user home directory") != nil {
			return "", err
		}

		// create hooks dir
		hookDir := filepath.Join(homeDir, hookFlag)
		return hookDir, nil
	}

	gitDir, err := getRepoRootDir()
	if handleError(err, "Failed to get repository root directory") != nil {
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
	if handleError(err, "Failed to get repository root directory with git command") != nil {
		return "", err
	}

	gitDir := filepath.Clean(byteOut.String())

	// remove /.git at last
	gitDir = filepath.Dir(gitDir)

	return gitDir, nil
}

func isHookExists(err error) bool {
	return errors.Is(err, errHooksExist)
}

func isConfExists(err error) bool {
	return errors.Is(err, errConfigExist)
}

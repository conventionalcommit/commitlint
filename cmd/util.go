package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

func runLint(confFilePath, fileInput string) (lintResult string, hasError bool, err error) {
	linter, format, err := getLinter(confFilePath)
	if err != nil {
		return "", false, err
	}

	commitMsg, err := getCommitMsg(fileInput)
	if err != nil {
		return "", false, err
	}

	res, err := linter.Lint(commitMsg)
	if err != nil {
		return "", false, err
	}

	resStr, err := format.Format(res)
	if err != nil {
		return "", false, err
	}
	return resStr, res.HasErrors(), nil
}

func getLinter(confFilePath string) (*lint.Linter, lint.Formatter, error) {
	conf, err := config.GetConfig(confFilePath)
	if err != nil {
		return nil, nil, err
	}

	format, err := config.GetFormatter(conf)
	if err != nil {
		return nil, nil, err
	}

	linter, err := config.GetLinter(conf)
	if err != nil {
		return nil, nil, err
	}

	return linter, format, nil
}

func setGitConf(hookDir string, isGlobal bool) error {
	var args = []string{"config"}
	if isGlobal {
		args = append(args, "--global")
	}
	args = append(args, "core.hooksPath", hookDir)

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getCommitMsg(fileInput string) (string, error) {
	commitMsg, err := readStdInPipe()
	if err != nil {
		return "", err
	}

	if commitMsg != "" {
		return commitMsg, nil
	}

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

func readStdInPipe() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	// user input from terminal
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// not handling this case
		return "", nil
	}

	// user input from stdin pipe
	readBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	s := string(readBytes)
	return strings.TrimSpace(s), nil
}

func getHookDir(isGlobal bool) (string, error) {
	baseDir := filepath.Clean(HookDir)

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

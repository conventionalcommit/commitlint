package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/lint"
)

const (
	// errExitCode represent error exit code
	errExitCode = 1
)

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

func runLint(confFilePath, fileInput string) (lintResult string, hasError bool, err error) {
	linter, format, err := getLinter(confFilePath)
	if err != nil {
		return "", false, err
	}

	commitMsg, err := getCommitMsg(fileInput)
	if err != nil {
		return "", false, err
	}

	result, err := linter.Lint(commitMsg)
	if err != nil {
		return "", false, err
	}

	output, err := format.Format(result)
	if err != nil {
		return "", false, err
	}
	return output, hasErrorSeverity(result), nil
}

func getLinter(confParam string) (*lint.Linter, lint.Formatter, error) {
	conf, err := getConfig(confParam)
	if err != nil {
		return nil, nil, err
	}

	format, err := config.GetFormatter(conf)
	if err != nil {
		return nil, nil, err
	}

	linter, err := config.NewLinter(conf)
	if err != nil {
		return nil, nil, err
	}

	return linter, format, nil
}

func getConfig(confParam string) (*lint.Config, error) {
	if confParam != "" {
		confParam = filepath.Clean(confParam)
		return config.Parse(confParam)
	}

	// If config param is empty, lookup for defaults
	conf, err := config.LookupAndParse()
	if err != nil {
		return nil, err
	}

	return conf, nil
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

	fileInput = filepath.Clean(fileInput)
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

func hasErrorSeverity(result *lint.Result) bool {
	for _, i := range result.Issues() {
		if i.Severity() == lint.SeverityError {
			return true
		}
	}
	return false
}

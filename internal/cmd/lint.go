package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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

	cleanMsg, err := cleanupMsg(commitMsg)
	if err != nil {
		return "", false, err
	}

	result, err := linter.ParseAndLint(cleanMsg)
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

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func cleanupMsg(dirtyMsg string) (string, error) {
	// commit msg cleanup in git is configurable: https://git-scm.com/docs/git-commit#Documentation/git-commit.txt---cleanupltmodegt
	// For now we do a combination of the "scissors" behavior and the "strip" behavior
	// * remove the scissors line and everything below
	// * strip leading and trailing empty lines
	// * strip commentary (lines stating with commentChar '#')
	// * strip trailing whitespace
	// * collapse consecutive empty lines
	// TODO: check via "git config --get" if any of those two hardcoded constants was reconfigured
	// TODO: find out if commit messages on windows actually

	gitCommentChar := "#"
	scissors := gitCommentChar + " ------------------------ >8 ------------------------"

	cleanMsg := ""
	lastLine := ""
	for _, line := range strings.Split(dirtyMsg, "\n") {
		if line == scissors {
			// remove everything below scissors (including the scissors line)
			break
		}
		if strings.HasPrefix(line, gitCommentChar) {
			// strip commentary
			continue
		}
		line = trimRightSpace(line)
		// strip trailing whitespace
		if lastLine == "" && line == "" {
			// strip leading empty lines
			// collapse consecutive empty lines
			continue
		}
		if cleanMsg == "" {
			cleanMsg = line
		} else {
			cleanMsg += "\n" + line
		}
		lastLine = line
	}
	if lastLine == "" {
		// strip trailing empty line
		cleanMsg = strings.TrimSuffix(cleanMsg, "\n")
	}
	return cleanMsg, nil
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

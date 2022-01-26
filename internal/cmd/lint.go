package cmd

import (
	"os"

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
	return resStr, hasErrorSeverity(res), nil
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

func hasErrorSeverity(res *lint.Failure) bool {
	for _, r := range res.Failures() {
		if r.Severity() == lint.SeverityError {
			return true
		}
	}
	return false
}

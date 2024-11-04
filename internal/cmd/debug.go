package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/conventionalcommit/commitlint/internal"
)

func printDebug() error {
	w := &strings.Builder{}
	w.WriteString("Commitlint Version: ")
	w.WriteString(internal.FullVersion())
	w.WriteByte('\n')

	gitVer, err := getGitVersion()
	if handleError(err, "Failed to get Git version") != nil {
		return err
	}

	localConf, err := getGitHookConfig(false)
	if handleError(err, "Failed to get local Git hook configuration") != nil {
		return err
	}

	globalConf, err := getGitHookConfig(true)
	if handleError(err, "Failed to get global Git hook configuration") != nil {
		return err
	}

	confFile, confType, err := internal.LookupConfigPath()
	if handleError(err, "Failed to lookup configuration path") != nil {
		return err
	}

	w.WriteString("Git Version: ")
	w.WriteString(gitVer)
	w.WriteByte('\n')

	w.WriteString("Local Hook: ")
	w.WriteString(localConf)
	w.WriteByte('\n')

	w.WriteString("Global Hook: ")
	w.WriteString(globalConf)

	switch confType {
	case internal.DefaultConfig:
		fmt.Fprintf(w, "\nConfig: Default")
	case internal.FileConfig:
		fmt.Fprintf(w, "\nConfig: %s - %s", confType, confFile)
	case internal.EnvConfig:
		fmt.Fprintf(w, "\nConfig: %s:%s - %s", confType, internal.CommitlintConfigEnv, confFile)
	}

	fmt.Println(w.String())
	return nil
}

func getGitVersion() (string, error) {
	b := &bytes.Buffer{}

	cmd := exec.Command("git", "version")
	cmd.Stdout = b
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if handleError(err, "Failed to execute 'git version' command") != nil {
		return "", err
	}

	ver := strings.ReplaceAll(b.String(), "git version ", "v")
	ver = strings.Trim(ver, "\n")
	return ver, nil
}

func getGitHookConfig(isGlobal bool) (string, error) {
	b := &bytes.Buffer{}

	var args = []string{"config"}
	if isGlobal {
		args = append(args, "--global")
	}
	args = append(args, "core.hooksPath")

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = b

	err := cmd.Run()
	if handleError(err, "Failed to execute 'git config core.hooksPath' command") != nil {
		return "", err
	}

	s := strings.TrimSpace(b.String())
	s = strings.Trim(s, "\n")

	return s, nil
}

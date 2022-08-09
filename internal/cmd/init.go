package cmd

import (
	"os"
	"os/exec"
)

// initLint is the callback function for init command
func initLint(confPath, hooksPath string, isGlobal, isReplace bool) error {
	hookDir, err := initHooks(confPath, hooksPath, isGlobal, isReplace)
	if err != nil {
		return err
	}
	return setGitConf(hookDir, isGlobal)
}

func setGitConf(hookDir string, isGlobal bool) error {
	args := []string{"config"}
	if isGlobal {
		args = append(args, "--global")
	}
	args = append(args, "core.hooksPath", hookDir)

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

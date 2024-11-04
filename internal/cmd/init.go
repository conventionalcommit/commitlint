package cmd

import (
	"os"
	"os/exec"
)

// initLint is the callback function for the init command
func initLint(confPath, hooksPath string, isGlobal, isReplace bool) error {
	hookDir, err := initHooks(confPath, hooksPath, isGlobal, isReplace)
	if handleError(err, "Failed to initialize hooks") != nil {
		return err
	}
	return handleError(setGitConf(hookDir, isGlobal), "Failed to set git configuration")
}

func setGitConf(hookDir string, isGlobal bool) error {
	args := []string{"config"}
	if isGlobal {
		args = append(args, "--global")
	}
	args = append(args, "core.hooksPath", hookDir)

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return handleError(cmd.Run(), "Failed to execute git config command")
}

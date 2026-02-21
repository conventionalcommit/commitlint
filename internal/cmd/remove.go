package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// removeLint is the callback function for the uninstall command.
func removeLint(isGlobal bool) error {
	return gitRemoveHook(isGlobal)
}

// gitRemoveHook removes commitlint from git config by unsetting core.hooksPath.
// It prompts for confirmation before making any changes.
// Hook files are left intact - the user can remove them manually if needed.
func gitRemoveHook(isGlobal bool) error {
	scope := "local"
	if isGlobal {
		scope = "global"
	}
	confirmed, err := promptConfirm(fmt.Sprintf("Unset core.hooksPath from %s git config?", scope))
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("aborted")
		return nil
	}

	if err := unsetGitConf(isGlobal); err != nil {
		return fmt.Errorf("could not unset git core.hooksPath: %w", err)
	}

	fmt.Println("commitlint hook removed successfully")
	fmt.Println("note: hook files were not removed - delete them manually if no longer needed")
	return nil
}

// unsetGitConf removes the core.hooksPath entry from git config (local or global).
func unsetGitConf(isGlobal bool) error {
	args := []string{"config"}
	if isGlobal {
		args = append(args, "--global")
	}
	args = append(args, "--unset", "core.hooksPath")

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// promptConfirm prints prompt and waits for the user to type y/yes or anything else.
// Returns true only when the user confirms with "y" or "yes" (case-insensitive).
func promptConfirm(prompt string) (bool, error) {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		// empty Enter is reported by fmt.Scanln as "unexpected newline"
		if err.Error() == "unexpected newline" {
			return false, nil
		}
		return false, err
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

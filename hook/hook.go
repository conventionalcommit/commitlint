// Package hook contains git hook related informations
package hook

import (
	"os"
	"path/filepath"
)

// CommitMsgHook represent commit-msg hook file name
const CommitMsgHook = "commit-msg"

// WriteToFile util func to write commit-msg hook to given file
func WriteToFile(hookDir string) (retErr error) {
	hookFilePath := filepath.Join(hookDir, filepath.Clean(CommitMsgHook))
	// if commit-msg already exists skip creating or overwriting it
	if _, err := os.Stat(hookFilePath); !os.IsNotExist(err) {
		return nil
	}
	// commit-msg needs to be executable
	file, err := os.OpenFile(hookFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer func() {
		err1 := file.Close()
		if retErr == nil && err1 != nil {
			retErr = err1
		}
	}()

	commitHook := commitMsgHook()
	_, err = file.WriteString(commitHook)
	return err
}

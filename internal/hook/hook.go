// Package hook contains git hook related informations
package hook

import (
	"os"
	"path/filepath"
)

// commitMsgHook represent commit-msg hook file name
const commitMsgHook = "commit-msg"

const hookMessage = `#!/bin/sh

if ! type commitlint >/dev/null 2>/dev/null; then
	echo ""
    echo "commitlint could not be found"
    echo "try again after installing commitlint or add commitlint to PATH"
	echo ""
    exit 2;
fi

commitlint lint --message $1

`

// WriteHooks write git hooks to the given outDir
func WriteHooks(outDir string) (retErr error) {
	hookFilePath := filepath.Join(outDir, commitMsgHook)
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

	_, err = file.WriteString(hookMessage)
	if err != nil {
		return err
	}
	return nil
}

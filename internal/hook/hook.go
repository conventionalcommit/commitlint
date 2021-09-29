// Package hook contains git hook related informations
package hook

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// commitMsgHook represent commit-msg hook file name
const commitMsgHook = "commit-msg"

// WriteHooks write git hooks to the given outDir
func WriteHooks(outDir, confPath string) (retErr error) {
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
	return writeTo(file, confPath)
}

// writeTo util func to write commit-msg hook to given io.Writer
func writeTo(wr io.Writer, confPath string) error {
	w := bufio.NewWriter(wr)

	w.WriteString("#!/bin/sh")
	w.WriteString("\n\ncommitlint lint")

	confPath = strings.TrimSpace(confPath)
	if confPath != "" {
		confPath = filepath.Clean(confPath)
		w.WriteString(` --config "` + confPath + `"`)
	}

	w.WriteString(" --message $1")
	w.WriteString("\n")

	return w.Flush()
}

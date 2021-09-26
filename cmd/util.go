package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"
)

func setGitConf(hookDir string, isGlobal bool) error {
	var args = []string{"config"}
	if isGlobal {
		args = append(args, "--global")
	}
	args = append(args, "core.hooksPath", hookDir)

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

func getRepoRootDir() (string, error) {
	byteOut := &bytes.Buffer{}

	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stdout = byteOut
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	gitDir := filepath.Clean(byteOut.String())

	// remove /.git at last
	gitDir = filepath.Dir(gitDir)

	return gitDir, nil
}

func formVersionInfo(version, commit, build string) string {
	versionTmpl := "%s - built from %s on %s"

	if build != "" {
		return fmt.Sprintf(versionTmpl, version, commit, build)
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Sprintf(versionTmpl, "master", "unknown", "unknown")
	}

	var commitInfo string
	if info.Main.Sum == "" {
		commitInfo = "(" + "checksum: unknown)"
	} else {
		commitInfo = "(" + "checksum: " + info.Main.Sum + ")"
	}

	return fmt.Sprintf(versionTmpl, info.Main.Version, commitInfo, "unknown")
}

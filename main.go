package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/conventionalcommit/commitlint/cmd"
)

// Build constants
// all variables are set during build
var (
	Version   string
	Commit    string
	BuildTime string
)

func init() {
	setVersionInfo()
}

func main() {
	app := cmd.New(Version, Commit, BuildTime)
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(cmd.ErrExitCode)
	}
}

func setVersionInfo() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		Version = info.Main.Version

		checkSum := "unknown"
		if info.Main.Sum != "" {
			checkSum = info.Main.Sum
		}

		Commit = "(" + "checksum: " + checkSum + ")"
	}

	if Version == "" {
		Version = "master"
	}

	if Commit == "" {
		Commit = "unknown"
	}

	if BuildTime == "" {
		BuildTime = "unknown"
	}
}

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
	if BuildTime == "" {
		info, ok := debug.ReadBuildInfo()
		if ok {
			checkSum := "unknown"
			if info.Main.Sum != "" {
				checkSum = info.Main.Sum
			}

			Version = info.Main.Version
			Commit = "(" + "checksum: " + checkSum + ")"
			BuildTime = "unknown"
			return
		}

		Version = "master"
		Commit = "unknown"
		BuildTime = "unknown"
	}
}

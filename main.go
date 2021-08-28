package main

import (
	"fmt"
	"os"

	"github.com/conventionalcommit/commitlint/cmd"
)

// Build constants
// all variables are set during build
var (
	Version   = "devel"
	Commit    = ""
	BuildTime = ""
)

func main() {
	app := cmd.New(Version, Commit, BuildTime)
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(cmd.ErrExitCode)
	}
}

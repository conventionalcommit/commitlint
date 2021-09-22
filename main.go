package main

import (
	"fmt"
	"os"

	"github.com/conventionalcommit/commitlint/cmd"
)

// Build constants
// all variables are set during build
var (
	Version   string
	Commit    string
	BuildTime string
)

func main() {
	app := cmd.NewWith(Version, Commit, BuildTime)
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(cmd.ErrExitCode)
	}
}

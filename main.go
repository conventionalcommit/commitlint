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

var verTmpl = `commitlint %s-%s %s
`

func main() {
	cmd.VersionCallback = func() error {
		fmt.Printf(verTmpl, Version, Commit, BuildTime)
		return nil
	}

	app := cmd.New()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(cmd.ErrExitCode)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/conventionalcommit/commitlint/cmd"
)

var errExitCode = 1

func main() {
	app := cmd.New()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(errExitCode)
	}
}

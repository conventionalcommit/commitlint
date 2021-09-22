// Package cmd contains commitlint cli
package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var versionTmpl = `commitlint version %s - built from %s on %s
`

// New returns commitlint cli.App
func New(versionNo, commitHash, builtTime string) *cli.App {
	versionInfo := fmt.Sprintf(versionTmpl, versionNo, commitHash, builtTime)

	cmds := []*cli.Command{
		createCmd(),
		initCmd(),
		lintCmd(),
		verifyCmd(),
		versionCmd(versionInfo),
	}

	app := &cli.App{
		Name:     "commitlint",
		Usage:    "linter for conventional commits",
		Commands: cmds,
		Action:   nil,
	}
	return app
}

func versionCmd(versionInfo string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "prints commitlint version",
		Action: func(c *cli.Context) error {
			fmt.Printf(versionInfo)
			return nil
		},
	}
}

func initCmd() *cli.Command {
	return &cli.Command{
		Name:   "init",
		Usage:  "setup commitlint for git repos",
		Action: initCallback,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "sets git hook in global config",
			},
		},
	}
}

func createCmd() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create commitlint config, hooks files",
		Subcommands: []*cli.Command{
			{
				Name:   "config",
				Usage:  "creates commitlint.yaml in current directory",
				Action: configCreateCallback,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "enabled",
						Aliases: []string{"e"},
						Usage:   "writes only default enabled rules to file",
						Value:   false,
					},
				},
			},
			{
				Name:   "hook",
				Usage:  "creates commit-msg file in current directory",
				Action: hookCreateCallback,
			},
		},
	}
}

func lintCmd() *cli.Command {
	return &cli.Command{
		Name:   "lint",
		Usage:  "lints commit message",
		Action: lintCallback,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c", "conf"},
				Value:   "",
				Usage:   "optional config file `conf.yaml`",
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m", "msg"},
				Value:   "",
				Usage:   "path to git commit message `FILE`",
			},
		},
	}
}

func verifyCmd() *cli.Command {
	return &cli.Command{
		Name:  "verify",
		Usage: "verifies commitlint config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c", "conf"},
				Value:   "",
				Usage:   "optional config file `conf.yaml`",
			},
		},
		Action: verifyCallback,
	}
}

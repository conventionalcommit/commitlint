package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint"
)

func getApp() *cli.App {
	createCmd := &cli.Command{
		Name:  "create",
		Usage: "for creating commitlint conf, hooks",
		Subcommands: []*cli.Command{
			{
				Name:  "config",
				Usage: "creates commitlint.yaml in current directory",
				Action: func(ctx *cli.Context) error {
					err := commitlint.DefaultConfToFile(defConfFileName)
					if err != nil {
						return cli.Exit(err, exitCode)
					}
					return nil
				},
			},
			{
				Name:  "hook",
				Usage: "creates commit-msg in current directory",
				Action: func(ctx *cli.Context) (retErr error) {
					err := commitlint.WriteHookToFile(commitMsgHook)
					if err != nil {
						return cli.Exit(err, exitCode)
					}
					return nil
				},
			},
		},
	}

	initCmd := &cli.Command{
		Name:   "init",
		Usage:  "setup commitlint for git repos",
		Action: initCallback,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "sets git hook config in global",
			},
		},
	}

	lintCmd := &cli.Command{
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

	versionCmd := &cli.Command{
		Name:  "version",
		Usage: "prints commitlint version",
		Action: func(c *cli.Context) error {
			fmt.Printf("commitlint - %s\n", commitlint.Version)
			return nil
		},
	}

	return &cli.App{
		Name:   "commitlint",
		Usage:  "linter for conventional commits",
		Action: nil,
		Commands: []*cli.Command{
			createCmd,
			initCmd,
			lintCmd,
			versionCmd,
		},
	}
}

package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint"
)

func getApp() *cli.App {
	createCmd := &cli.Command{
		Name:  "create",
		Usage: "helpers for initializing commitlint",
		Subcommands: []*cli.Command{
			{
				Name:   "config",
				Usage:  "creates default commitlint.yaml in current directory",
				Action: initConfCallback,
			},
			{
				Name:   "hook",
				Usage:  "creates commit hook and sets git config",
				Action: initHookCallback,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "global",
						Aliases: []string{"g"},
						Usage:   "sets git hook config in global`",
					},
				},
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
		Name:     "commitlint",
		Usage:    "linter for conventional commits",
		Action:   nil,
		Commands: []*cli.Command{lintCmd, createCmd, versionCmd},
	}
}

package cmd

import (
	"fmt"
	"os"

	"github.com/conventionalcommit/commitlint/internal"
	cli "github.com/urfave/cli/v2"
)

// Run runs commitlint cli with os.Args
func Run() error {
	return newCliApp().Run(os.Args)
}

// newCliApp returns commitlint cli.App
func newCliApp() *cli.App {
	cmds := []*cli.Command{
		newInitCmd(),
		newRemoveCmd(),
		newLintCmd(),
		newConfigCmd(),
		newHookCmd(),
		newDebugCmd(),
	}

	app := &cli.App{
		Name:     "commitlint",
		Usage:    "Lint commit messages using Conventional Commits rules",
		Commands: cmds,
		Version:  internal.FullVersion(),
	}
	return app
}

func newLintCmd() *cli.Command {
	return &cli.Command{
		Name:        "lint",
		Usage:       "Check a commit message",
		Description: "Reads from stdin (piped), or --message file, or .git/COMMIT_EDITMSG (in that order).",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Use this config `FILE` instead of auto-detected one",
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m", "msg"},
				Usage:   "Read commit message from `FILE`",
			},
		},
		Action: func(ctx *cli.Context) error {
			confFilePath := ctx.String("config")
			fileInput := ctx.String("message")
			return lintMsg(confFilePath, fileInput)
		},
	}
}

func newInitCmd() *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "Set up commitlint for a git repository",
		Description: "Creates the commit-msg hook and points git to it.\nUse --global to apply across all your repositories.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "Set up for all repositories (uses global git config)",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Pass a config `FILE` to the hook",
			},
			&cli.BoolFlag{
				Name:    "replace",
				Aliases: []string{"r"},
				Usage:   "Overwrite existing hook files",
			},
			&cli.StringFlag{
				Name:    "hookspath",
				Aliases: []string{"p"},
				Usage:   "Where to write hook files (default: .commitlint/hooks)",
			},
		},
		Action: func(ctx *cli.Context) error {
			confPath := ctx.String("config")
			isGlobal := ctx.Bool("global")
			isReplace := ctx.Bool("replace")
			hooksPath := ctx.String("hookspath")

			err := initLint(confPath, hooksPath, isGlobal, isReplace)
			if err != nil {
				if isHookExists(err) {
					fmt.Println("commitlint init failed: hook files already exist")
					fmt.Println("use --replace to overwrite them")
					return nil
				}
				return err
			}

			fmt.Println("commitlint init successfully")
			return nil
		},
	}
}

func newConfigCmd() *cli.Command {
	createCmd := &cli.Command{
		Name:  "create",
		Usage: "Create a default config file",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "replace",
				Aliases: []string{"r"},
				Usage:   "Overwrite if the file already exists",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "Output file name",
				Value: ".commitlint.yaml",
			},
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Write all settings (not just enabled rules)",
			},
		},
		Action: func(ctx *cli.Context) error {
			isReplace := ctx.Bool("replace")
			fileName := ctx.String("file")
			all := ctx.Bool("all")
			err := configCreate(fileName, isReplace, all)
			if err != nil {
				if isConfExists(err) {
					fmt.Println("config file already exists")
					fmt.Println("use --replace to overwrite it")
					return nil
				}
				return err
			}
			fmt.Println("config file created")
			return nil
		},
	}

	checkCmd := &cli.Command{
		Name:      "check",
		Usage:     "Check if a config file is valid",
		ArgsUsage: "<config-file>",
		Action: func(ctx *cli.Context) error {
			confFile := ctx.Args().First()
			if confFile == "" {
				return fmt.Errorf("please provide a config file path\n\nUsage: commitlint config check <config-file>")
			}
			errs := configCheck(confFile)
			if len(errs) == 0 {
				fmt.Printf("%s: valid\n", confFile)
				return nil
			}
			if len(errs) == 1 {
				return errs[0]
			}
			merr := multiError(errs)
			return &merr
		},
	}

	return &cli.Command{
		Name:        "config",
		Usage:       "Manage configuration",
		Subcommands: []*cli.Command{createCmd, checkCmd},
	}
}

func newHookCmd() *cli.Command {
	return &cli.Command{
		Name:  "hook",
		Usage: "Create the commit-msg git hook",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "replace",
				Aliases: []string{"r"},
				Usage:   "Overwrite existing hook files",
			},
			&cli.StringFlag{
				Name:    "hookspath",
				Aliases: []string{"p"},
				Usage:   "Where to write hook files (default: .commitlint/hooks)",
			},
		},
		Action: func(ctx *cli.Context) error {
			isReplace := ctx.Bool("replace")
			hooksPath := ctx.String("hookspath")
			err := hookCreate(hooksPath, isReplace)
			if err != nil {
				if isHookExists(err) {
					fmt.Println("hook files already exist")
					fmt.Println("use --replace to overwrite them")
					return nil
				}
				return err
			}
			fmt.Println("hooks created")
			return nil
		},
	}
}

func newRemoveCmd() *cli.Command {
	return &cli.Command{
		Name:        "remove",
		Usage:       "Remove commitlint from git config",
		Description: "Unset git's core.hooksPath so commits are no longer linted.\nHook files are left intact.\nUse --global to remove globally configured hooks.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "Remove from global git config",
			},
		},
		Action: func(ctx *cli.Context) error {
			isGlobal := ctx.Bool("global")
			return removeLint(isGlobal)
		},
	}
}

func newDebugCmd() *cli.Command {
	return &cli.Command{
		Name:  "debug",
		Usage: "Show debug info (version, hooks, config)",
		Action: func(ctx *cli.Context) error {
			return printDebug()
		},
	}
}

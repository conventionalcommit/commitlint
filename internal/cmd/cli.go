// Package cmd contains commitlint cli
package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint/internal"
)

// newCliApp returns commitlint cli.App
func newCliApp() *cli.App {
	cmds := []*cli.Command{
		newInitCmd(),
		newLintCmd(),
		newConfigCmd(),
		newHookCmd(),
		newDebugCmd(),
	}

	app := &cli.App{
		Name:     "commitlint",
		Usage:    "linter for conventional commits",
		Commands: cmds,
		Version:  internal.FullVersion(),
	}
	return app
}

func newLintCmd() *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "Check commit message against lint rules",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "",
				Usage:   "optional config file `conf.yaml`",
			},
			&cli.StringFlag{
				Name:    "message",
				Aliases: []string{"m", "msg"},
				Value:   "",
				Usage:   "path to commit message `FILE`",
			},
		},
		Action: func(ctx *cli.Context) error {
			confFilePath := ctx.String("config")
			fileInput := ctx.String("message")
			err := lintMsg(confFilePath, fileInput)
			return handleError(err, "Failed to run lint command")
		},
	}
}

func newInitCmd() *cli.Command {
	confFlag := newConfFlag()
	replaceFlag := newReplaceFlag()
	hooksFlag := newHooksPathFlag()

	globalFlag := &cli.BoolFlag{
		Name:    "global",
		Aliases: []string{"g"},
		Usage:   "Sets git hook in global config",
	}

	return &cli.Command{
		Name:  "init",
		Usage: "Setup commitlint for git repos",
		Flags: []cli.Flag{globalFlag, confFlag, replaceFlag, hooksFlag},
		Action: func(ctx *cli.Context) error {
			confPath := ctx.String("config")
			isGlobal := ctx.Bool("global")
			isReplace := ctx.Bool("replace")
			hooksPath := ctx.String("hookspath")

			err := initLint(confPath, hooksPath, isGlobal, isReplace)
			if handleError(err, "Failed to initialize commitlint") != nil {
				if isHookExists(err) {
					fmt.Println("commitlint init failed")
					fmt.Println("run with --replace to replace existing files")
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
		Usage: "Creates default config in current directory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "replace",
				Aliases: []string{"r"},
				Usage:   "Replace conf file if already exists",
				Value:   false,
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "Config file name",
				Value: ".commitlint.yaml",
			},
		},
		Action: func(ctx *cli.Context) error {
			isReplace := ctx.Bool("replace")
			fileName := ctx.String("file")
			err := configCreate(fileName, isReplace)
			if handleError(err, "Failed to create config file") != nil {
				if isConfExists(err) {
					fmt.Println("config create failed")
					fmt.Println("run with --replace to replace existing file")
					return nil
				}
				return err
			}
			fmt.Println("config file created")
			return nil
		},
	}

	checkCmd := &cli.Command{
		Name:  "check",
		Usage: "Checks if given config is valid",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "config file `conf.yaml`",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			confFile := ctx.String("config")
			errs := configCheck(confFile)
			if len(errs) == 0 {
				fmt.Printf("%s config is valid\n", confFile)
				return nil
			}
			if len(errs) == 1 {
				return handleError(errs[0], "Config check failed")
			}
			merr := multiError(errs)
			return handleError(&merr, "Config check failed")
		},
	}

	return &cli.Command{
		Name:        "config",
		Usage:       "Manage commitlint config",
		Subcommands: []*cli.Command{createCmd, checkCmd},
	}
}

func newHookCmd() *cli.Command {
	replaceFlag := newReplaceFlag()
	hooksFlag := newHooksPathFlag()

	createCmd := &cli.Command{
		Name:  "create",
		Usage: "Creates git hook files in current directory",
		Flags: []cli.Flag{replaceFlag, hooksFlag},
		Action: func(ctx *cli.Context) error {
			isReplace := ctx.Bool("replace")
			hooksPath := ctx.String("hookspath")
			err := hookCreate(hooksPath, isReplace)
			if handleError(err, "Failed to create hooks") != nil {
				if isHookExists(err) {
					fmt.Println("create failed. hook files already exist")
					fmt.Println("run with --replace to replace existing hook files")
					return nil
				}
				return err
			}
			fmt.Println("hooks created")
			return nil
		},
	}

	return &cli.Command{
		Name:        "hook",
		Usage:       "Manage commitlint git hooks",
		Subcommands: []*cli.Command{createCmd},
	}
}

func newDebugCmd() *cli.Command {
	return &cli.Command{
		Name:  "debug",
		Usage: "prints useful information for debugging",
		Action: func(ctx *cli.Context) error {
			return handleError(printDebug(), "Debugging information failed")
		},
	}
}

func newConfFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "",
		Usage:   "Optional config file `conf.yaml` which will be passed to 'commitlint lint'. Check config precedence",
	}
}

func newHooksPathFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "hookspath",
		Value: "",
		Usage: "Optional hookspath to install git hooks",
	}
}

func newReplaceFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "replace",
		Usage: "Replace hook files if already exists",
	}
}

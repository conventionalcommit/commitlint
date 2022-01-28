// Package cmd contains commitlint cli
package cmd

import (
	"fmt"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/conventionalcommit/commitlint/internal"
)

// Run runs commitlint cli with os.Args
func Run() error {
	return newCmd().Run(os.Args)
}

// newCmd returns commitlint cli.App
func newCmd() *cli.App {
	cmds := []*cli.Command{
		initCmd(),
		lintCmd(),
		configCmd(),
		hookCmd(),
	}

	app := &cli.App{
		Name:     "commitlint",
		Usage:    "linter for conventional commits",
		Commands: cmds,
		Version:  internal.FullVersion(),
	}
	return app
}

func lintCmd() *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "Check commit message against lint rules",
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
				Usage:   "path to commit message `FILE`",
			},
		},
		Action: func(ctx *cli.Context) error {
			confFilePath := ctx.String("config")
			fileInput := ctx.String("message")
			return lintMsg(confFilePath, fileInput)
		},
	}
}

func initCmd() *cli.Command {
	confFlag := formConfFlag()
	replaceFlag := formReplaceFlag()

	globalFlag := &cli.BoolFlag{
		Name:    "global",
		Aliases: []string{"g"},
		Usage:   "Sets git hook in global config",
	}

	return &cli.Command{
		Name:  "init",
		Usage: "Setup commitlint for git repos",
		Flags: []cli.Flag{globalFlag, confFlag, replaceFlag},
		Action: func(ctx *cli.Context) error {
			confPath := ctx.String("config")
			isGlobal := ctx.Bool("global")
			isReplace := ctx.Bool("replace")

			err := initLint(confPath, isGlobal, isReplace)
			if err != nil {
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

func configCmd() *cli.Command {
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
		},
		Action: func(ctx *cli.Context) error {
			isReplace := ctx.Bool("replace")
			err := configCreate(isReplace)
			if err != nil {
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
				Aliases:  []string{"c", "conf"},
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
				return errs[0]
			}
			merr := multiError(errs)
			return &merr
		},
	}

	return &cli.Command{
		Name:        "config",
		Usage:       "Manage commitlint config",
		Subcommands: []*cli.Command{createCmd, checkCmd},
	}
}

func hookCmd() *cli.Command {
	confFlag := formConfFlag()
	replaceFlag := formReplaceFlag()

	createCmd := &cli.Command{
		Name:  "create",
		Usage: "Creates git hook files in current directory",
		Flags: []cli.Flag{confFlag, replaceFlag},
		Action: func(ctx *cli.Context) error {
			confPath := ctx.String("config")
			isReplace := ctx.Bool("replace")
			err := hookCreate(confPath, isReplace)
			if err != nil {
				if isHookExists(err) {
					fmt.Println("create failed. hook files already exists")
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

func formConfFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c", "conf"},
		Value:   "",
		Usage:   "Optional config file `conf.yaml` which will be passed to 'commitlint lint'. Check config precedence",
	}
}

func formReplaceFlag() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "replace",
		Usage: "Replace hook files if already exists",
	}
}

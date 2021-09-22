// Package cmd contains commitlint cli
package cmd

import (
	"github.com/urfave/cli/v2"
)

// New returns commitlint cli.App
func New() *cli.App {
	return NewWith("", "", "")
}

// NewWith returns commitlint cli.App with version info
func NewWith(versionNo, commitHash, builtTime string) *cli.App {
	versionInfo := formVersionInfo(versionNo, commitHash, builtTime)

	cmds := []*cli.Command{
		initCmd(),
		lintCmd(),
		createCmd(),
		verifyCmd(),
	}

	app := &cli.App{
		Name:     "commitlint",
		Usage:    "linter for conventional commits",
		Commands: cmds,
		Version:  versionInfo,
	}
	return app
}

func initCmd() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "setup commitlint for git repos",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "global",
				Aliases: []string{"g"},
				Usage:   "sets git hook in global config",
			},
		},
		Action: func(ctx *cli.Context) error {
			isGlobal := ctx.Bool("global")
			return Init(isGlobal)
		},
	}
}

func createCmd() *cli.Command {
	configCmd := &cli.Command{
		Name:  "config",
		Usage: "creates commitlint.yaml in current directory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "enabled",
				Aliases: []string{"e"},
				Usage:   "writes only default enabled rules to file",
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			isOnlyEnabled := ctx.Bool("enabled")
			return CreateConfig(isOnlyEnabled)
		},
	}

	hookCmd := &cli.Command{
		Name:  "hook",
		Usage: "creates commit-msg file in current directory",
		Action: func(ctx *cli.Context) error {
			return CreateHook()
		},
	}

	return &cli.Command{
		Name:  "create",
		Usage: "create commitlint config, hooks files",
		Subcommands: []*cli.Command{
			configCmd,
			hookCmd,
		},
	}
}

func lintCmd() *cli.Command {
	return &cli.Command{
		Name:  "lint",
		Usage: "lints commit message",
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
		Action: func(ctx *cli.Context) error {
			confFilePath := ctx.String("config")
			fileInput := ctx.String("message")
			return Lint(confFilePath, fileInput)
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
		Action: func(ctx *cli.Context) error {
			confFilePath := ctx.String("config")
			return VerifyConfig(confFilePath)
		},
	}
}

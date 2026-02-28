package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/conventionalcommit/commitlint/changelog"
	"github.com/conventionalcommit/commitlint/config"
	"github.com/conventionalcommit/commitlint/registry"
	cli "github.com/urfave/cli/v2"
)

func newChangelogCmd() *cli.Command {
	return &cli.Command{
		Name:        "changelog",
		Usage:       "Generate changelog from conventional commits",
		Description: "Generates a changelog grouped by commit type.\nSmart mode: if no flags, generates full changelog for all versions.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Config `FILE` path (default: auto-detect)",
			},
			&cli.StringFlag{
				Name:  "from",
				Usage: "Start reference (tag/commit/branch, exclusive)",
			},
			&cli.StringFlag{
				Name:  "to",
				Usage: "End reference (default: HEAD)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output `FILE` (default: stdout)",
			},
			&cli.StringFlag{
				Name:  "format",
				Usage: "Output format: markdown, json (default: from config)",
			},
		},
		Action: func(ctx *cli.Context) error {
			confPath := ctx.String("config")
			from := ctx.String("from")
			to := ctx.String("to")
			output := ctx.String("output")
			format := ctx.String("format")

			return runChangelog(confPath, from, to, output, format)
		},
	}
}

func runChangelog(confPath, from, to, output, format string) error {
	// load changelog config
	clConf, err := loadChangelogConfig(confPath)
	if err != nil {
		return err
	}

	// override format if specified via CLI
	if format != "" {
		clConf.Formatter = format
	}

	// get current directory as repo dir
	repoDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// create generator
	gen, err := changelog.New(repoDir, clConf)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// generate changelog
	cl, err := gen.GenerateSmart(from, to)
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	// get formatter
	f, ok := registry.GetChangelogFormatter(clConf.Formatter)
	if !ok {
		return fmt.Errorf("unknown changelog formatter '%s'", clConf.Formatter)
	}

	// format output
	result, err := f.Format(cl)
	if err != nil {
		return fmt.Errorf("failed to format changelog: %w", err)
	}

	// write output
	if output != "" && output != "-" {
		output = filepath.Clean(output)
		err = os.WriteFile(output, []byte(result), 0o644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("changelog written to %s\n", output)
		return nil
	}

	fmt.Print(result)
	return nil
}

func loadChangelogConfig(confPath string) (*changelog.Config, error) {
	if confPath != "" {
		confPath = filepath.Clean(confPath)
		conf, err := config.Parse(confPath)
		if err != nil {
			return nil, err
		}
		return conf.Changelog, nil
	}

	// Try to get config from the normal lookup
	conf, err := config.LookupAndParse()
	if err != nil {
		return nil, err
	}

	return conf.Changelog, nil
}

func printChangelogDebug() error {
	fmt.Println()
	fmt.Println("Changelog:")

	// load config
	conf, err := config.LookupAndParse()
	if err != nil {
		fmt.Printf("  Config Error: %s\n", err)
		return nil
	}

	clConf := conf.Changelog

	fmt.Printf("  Formatter:    %s\n", clConf.Formatter)

	// git info
	repoDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("  Git:          error: %s\n", err)
		return nil
	}

	gen, err := changelog.New(repoDir, clConf)
	if err != nil {
		fmt.Printf("  Git:          error: %s\n", err)
		return nil
	}

	info, err := gen.Debug()
	if err != nil {
		fmt.Printf("  Git:          error: %s\n", err)
		return nil
	}

	fmt.Printf("  Remote URL:   %s\n", valueOrNA(info.RemoteURL))
	fmt.Printf("  Host Type:    %s\n", string(info.HostType))
	fmt.Printf("  Commit URL:   %s\n", valueOrNA(info.CommitURL))
	fmt.Printf("  Compare URL:  %s\n", valueOrNA(info.CompareURL))
	fmt.Printf("  HEAD:         %s\n", valueOrNA(info.Head))
	fmt.Printf("  Latest Tag:   %s\n", valueOrNA(info.LatestTag))
	fmt.Printf("  Total Tags:   %d\n", info.TotalTags)

	// types
	fmt.Println("  Types:")
	for _, t := range info.Types {
		hidden := ""
		if t.Hidden {
			hidden = " (hidden)"
		}
		fmt.Printf("    %s → %s%s\n", t.Type, t.Header, hidden)
	}

	// issue prefixes
	if len(clConf.IssuePrefixes) > 0 {
		fmt.Printf("  Issue Prefixes: %s\n", strings.Join(clConf.IssuePrefixes, ", "))
	}

	return nil
}

func valueOrNA(s string) string {
	if s == "" {
		return "N/A"
	}
	return s
}

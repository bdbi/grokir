package main

import (
	"flag"
	"fmt"
	"os"

	"grokir/internal/cli/command"
	_ "grokir/internal/cli/commands"
	"grokir/internal/grokipedia"
)

var (
	version = "dev"
	date    = "unknown"
)

func usage() {
	fmt.Fprintf(os.Stderr, `grokir - CLI for Grokipedia
Version: %s
Build date: %s

Usage:
  grokir search <query> [-l <num>] [-o <num>] [--json]
  grokir page <slug> [--json]
  grokir version

Commands:
  search   Search articles on Grokipedia
  page     Show a page by slug
  version  Show version and build date

Options:
  -l       Maximum number of results (default: 10)
  -o       Offset for pagination (default: 0)
  --json   JSON output
`, version, date)
}

func main() {
	flag.Usage = usage

	jsonOutput := flag.Bool("json", false, "JSON output")

	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}

	cmd := flag.Arg(0)
	args := flag.Args()[1:]

	if cmd == "version" {
		fmt.Printf("Version: %s\nBuild date: %s\n", version, date)
		return
	}

	client := grokipedia.NewClient()

	var outputMode command.OutputMode
	if *jsonOutput {
		outputMode = command.OutputJSON
	} else {
		outputMode = command.OutputText
	}

	rt := command.Runtime{
		Client: client,
		Output: outputMode,
	}

	err := command.Run(rt, cmd, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if cmdErr, ok := err.(*command.Error); ok && cmdErr.IsUsage() {
			usage()
		}
		os.Exit(1)
	}
}

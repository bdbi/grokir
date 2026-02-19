package commands

import (
	"flag"
	"fmt"
	"strings"

	"grokir/internal/cli/command"
	"grokir/internal/cli/formatter"
)

type searchCommand struct{}

func init() {
	command.Register(&searchCommand{})
}

func (c *searchCommand) Name() string {
	return "search"
}

func (c *searchCommand) Usage() string {
	return "grokir search <query> [-l <num>] [-o <num>]"
}

func (c *searchCommand) Run(rt command.Runtime, args []string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Print(c.Usage())
	}
	limit := fs.Int("l", 10, "maximum number of results")
	offset := fs.Int("o", 0, "offset for pagination")

	if err := fs.Parse(args); err != nil {
		return command.NewUsageError(err.Error())
	}

	if fs.NArg() < 1 {
		return command.NewUsageError("missing search query")
	}

	query := strings.Join(fs.Args(), " ")

	results, err := rt.Client.Search(query, *limit, *offset)
	if err != nil {
		return command.NewRuntimeError("search error: %v", err)
	}

	f := formatter.NewSearchFormatter(rt.Output)

	if len(results) == 0 {
		fmt.Print(f.NoResults())
		return nil
	}

	output, err := f.FormatSearch(results)
	if err != nil {
		return command.NewRuntimeError("formatting error: %v", err)
	}

	fmt.Print(output)
	return nil
}

var _ command.Command = (*searchCommand)(nil)

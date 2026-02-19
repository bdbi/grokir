package commands

import (
	"fmt"

	"grokir/internal/cli/command"
	"grokir/internal/cli/formatter"
)

type pageCommand struct{}

func init() {
	command.Register(&pageCommand{})
}

func (c *pageCommand) Name() string {
	return "page"
}

func (c *pageCommand) Usage() string {
	return "grokir page <slug>"
}

func (c *pageCommand) Run(rt command.Runtime, args []string) error {
	if len(args) < 1 {
		return command.NewUsageError("missing page slug")
	}

	slug := args[0]

	page, err := rt.Client.GetPage(slug, true, false)
	if err != nil {
		return command.NewRuntimeError("page retrieval error: %v", err)
	}

	f := formatter.NewPageFormatter(rt.Output)

	output, err := f.FormatPage(page)
	if err != nil {
		return command.NewRuntimeError("formatting error: %v", err)
	}

	fmt.Print(output)
	return nil
}

var _ command.Command = (*pageCommand)(nil)

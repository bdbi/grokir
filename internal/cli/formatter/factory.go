package formatter

import "grokir/internal/cli/command"

// NewSearchFormatter returns a SearchFormatter for the given output mode.
func NewSearchFormatter(mode command.OutputMode) command.SearchFormatter {
	switch mode {
	case command.OutputJSON:
		return NewJSON()
	default:
		return NewText()
	}
}

// NewPageFormatter returns a PageFormatter for the given output mode.
func NewPageFormatter(mode command.OutputMode) command.PageFormatter {
	switch mode {
	case command.OutputJSON:
		return NewJSON()
	default:
		return NewText()
	}
}

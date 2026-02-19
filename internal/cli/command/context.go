package command

import (
	"grokir/internal/grokipedia"
)

// OutputMode represents the output format for command results.
type OutputMode string

const (
	OutputText OutputMode = "text"
	OutputJSON OutputMode = "json"
)

// SearchFormatter defines the interface for formatting search results.
type SearchFormatter interface {
	FormatSearch([]grokipedia.SearchResult) (string, error)
	NoResults() string
}

// PageFormatter defines the interface for formatting page content.
type PageFormatter interface {
	FormatPage(*grokipedia.Page) (string, error)
}

// Runtime holds the shared dependencies required by all commands.
type Runtime struct {
	Client *grokipedia.Client
	Output OutputMode
}

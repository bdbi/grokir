package formatter

import (
	"encoding/json"
	"fmt"

	"grokir/internal/grokipedia"
)

// JSONFormatter formats output as indented JSON.
type JSONFormatter struct{}

func NewJSON() *JSONFormatter {
	return &JSONFormatter{}
}

// FormatSearch renders search results as a JSON array.
func (f *JSONFormatter) FormatSearch(results []grokipedia.SearchResult) (string, error) {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON error: %w", err)
	}
	return string(data), nil
}

// FormatPage renders a page as a JSON object.
func (f *JSONFormatter) FormatPage(page *grokipedia.Page) (string, error) {
	data, err := json.MarshalIndent(page, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON error: %w", err)
	}
	return string(data), nil
}

// NoResults returns an empty JSON array for when search returns no results.
func (f *JSONFormatter) NoResults() string {
	return "[]"
}

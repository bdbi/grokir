package formatter

import (
	"fmt"
	"strings"

	"grokir/internal/grokipedia"
)

// TextFormatter formats output as human-readable plain text.
type TextFormatter struct{}

func NewText() *TextFormatter {
	return &TextFormatter{}
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func truncate(s string, maxLen int) string {
	s = normalize(s)
	if len(s) <= maxLen {
		return s
	}
	trunc := strings.TrimRight(s[:maxLen], " ")
	return trunc + "..."
}

func formatViews(n int64) string {
	s := fmt.Sprintf("%d", n)
	var result []byte
	length := len(s)
	for i := 0; i < length; i++ {
		if i > 0 && (length-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, s[i])
	}
	return string(result)
}

// FormatSearch renders search results as a compact terminal-friendly list.
func (f *TextFormatter) FormatSearch(results []grokipedia.SearchResult) (string, error) {
	var b strings.Builder
	for i, r := range results {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf("%d) %s\n", i+1, r.Title))
		b.WriteString(fmt.Sprintf("   slug: %s | relevance: %.2f | views: %s\n",
			r.Slug, r.RelevanceScore, formatViews(r.ViewCount)))
		if r.Snippet != "" {
			b.WriteString(fmt.Sprintf("   %s\n", truncate(r.Snippet, 200)))
		}
	}
	return b.String(), nil
}

// FormatPage renders a page as a readable text document.
func (f *TextFormatter) FormatPage(page *grokipedia.Page) (string, error) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Title: %s\n", page.Title))
	b.WriteString(fmt.Sprintf("Slug: %s\n", page.Slug))
	if page.Description != "" {
		b.WriteString(fmt.Sprintf("\n%s\n", normalize(page.Description)))
	}
	b.WriteString(strings.Repeat("-", 40))
	b.WriteString("\n")
	if page.Content != "" {
		b.WriteString(page.Content)
	}
	return b.String(), nil
}

// NoResults returns a human-readable message when search returns no results.
func (f *TextFormatter) NoResults() string {
	return "No results found. Try different keywords.\n"
}

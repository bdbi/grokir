package formatter

import (
	"encoding/json"
	"testing"

	"grokir/internal/grokipedia"
)

func TestJSONFormatter_FormatSearch(t *testing.T) {
	f := NewJSON()

	tests := []struct {
		name    string
		results []grokipedia.SearchResult
	}{
		{
			name: "single result",
			results: []grokipedia.SearchResult{
				{
					Slug:           "test",
					Title:          "Test",
					Snippet:        "A snippet",
					RelevanceScore: 0.9,
					ViewCount:      100,
				},
			},
		},
		{
			name:    "empty results",
			results: []grokipedia.SearchResult{},
		},
		{
			name: "multiple results",
			results: []grokipedia.SearchResult{
				{Slug: "one", Title: "One", RelevanceScore: 0.8, ViewCount: 50},
				{Slug: "two", Title: "Two", RelevanceScore: 0.7, ViewCount: 75},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.FormatSearch(tt.results)
			if err != nil {
				t.Fatalf("FormatSearch() error = %v", err)
			}

			var parsed []grokipedia.SearchResult
			if err := json.Unmarshal([]byte(got), &parsed); err != nil {
				t.Fatalf("FormatSearch() returned invalid JSON: %v", err)
			}

			if len(parsed) != len(tt.results) {
				t.Errorf("FormatSearch() returned %d items, want %d", len(parsed), len(tt.results))
			}
		})
	}
}

func TestJSONFormatter_FormatPage(t *testing.T) {
	f := NewJSON()

	tests := []struct {
		name string
		page *grokipedia.Page
	}{
		{
			name: "full page",
			page: &grokipedia.Page{
				Title:       "Kubernetes",
				Slug:        "kubernetes",
				Description: "Container orchestrator",
				Content:     "Content here",
			},
		},
		{
			name: "minimal page",
			page: &grokipedia.Page{
				Title: "Minimal",
				Slug:  "minimal",
			},
		},
		{
			name: "nil page",
			page: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.FormatPage(tt.page)
			if err != nil {
				t.Fatalf("FormatPage() error = %v", err)
			}

			if tt.page == nil {
				if got != "null" {
					t.Fatalf("FormatPage(nil) = %q, want %q", got, "null")
				}
				return
			}

			var parsed grokipedia.Page
			if err := json.Unmarshal([]byte(got), &parsed); err != nil {
				t.Fatalf("FormatPage() returned invalid JSON: %v", err)
			}

			if parsed.Title != tt.page.Title {
				t.Errorf("Title = %q, want %q", parsed.Title, tt.page.Title)
			}
			if parsed.Slug != tt.page.Slug {
				t.Errorf("Slug = %q, want %q", parsed.Slug, tt.page.Slug)
			}
		})
	}
}

func TestJSONFormatter_NoResults(t *testing.T) {
	f := NewJSON()
	got := f.NoResults()
	want := "[]"
	if got != want {
		t.Errorf("NoResults() = %q, want %q", got, want)
	}
}

package formatter

import (
	"strings"
	"testing"

	"grokir/internal/grokipedia"
)

func TestTextFormatter_FormatSearch(t *testing.T) {
	f := NewText()

	tests := []struct {
		name    string
		results []grokipedia.SearchResult
		wantSub []string
		notWant []string
	}{
		{
			name: "single result with snippet",
			results: []grokipedia.SearchResult{
				{
					Slug:           "kubernetes-scheduling",
					Title:          "Kubernetes Scheduling",
					Snippet:        "The scheduler decides where pods run by evaluating resource requests and limits.",
					RelevanceScore: 0.95,
					ViewCount:      12500,
				},
			},
			wantSub: []string{
				"1) Kubernetes Scheduling",
				"slug: kubernetes-scheduling",
				"relevance: 0.95",
				"views: 12,500",
			},
			notWant: []string{"\n\n"},
		},
		{
			name: "multiple results",
			results: []grokipedia.SearchResult{
				{Title: "First", Slug: "first", Snippet: "Snippet one", RelevanceScore: 0.9, ViewCount: 100},
				{Title: "Second", Slug: "second", Snippet: "Snippet two", RelevanceScore: 0.8, ViewCount: 200},
			},
			wantSub: []string{
				"1) First", "2) Second",
			},
		},
		{
			name: "result without snippet",
			results: []grokipedia.SearchResult{
				{Title: "No Snippet", Slug: "no-snippet", Snippet: "", RelevanceScore: 0.5, ViewCount: 50},
			},
			wantSub: []string{"1) No Snippet"},
			notWant: []string{"   \n"},
		},
		{
			name: "long snippet is truncated",
			results: []grokipedia.SearchResult{
				{
					Title: "Long Snippet",
					Slug:  "long-snippet",
					Snippet: "This is a long snippet that is intentionally larger than two hundred characters so we can verify truncation behavior in the text formatter. " +
						"It should end with an ellipsis and should not include the entire original text after the cutoff point.",
				},
			},
			wantSub: []string{"..."},
			notWant: []string{"entire original text after the cutoff point."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.FormatSearch(tt.results)
			if err != nil {
				t.Fatalf("FormatSearch() error = %v", err)
			}

			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("FormatSearch() missing substring %q in:\n%s", sub, got)
				}
			}

			for _, sub := range tt.notWant {
				if strings.Contains(got, sub) {
					t.Errorf("FormatSearch() should not contain %q in:\n%s", sub, got)
				}
			}
		})
	}
}

func TestTextFormatter_FormatPage(t *testing.T) {
	f := NewText()

	tests := []struct {
		name       string
		page       *grokipedia.Page
		wantSub    []string
		notWantSub []string
	}{
		{
			name: "full page with description",
			page: &grokipedia.Page{
				Title:       "Kubernetes",
				Slug:        "kubernetes",
				Description: "Container orchestration system.",
				Content:     "Content goes here.",
			},
			wantSub: []string{
				"Title: Kubernetes",
				"Slug: kubernetes",
				"Container orchestration system.",
				"Content goes here.",
			},
		},
		{
			name: "page without description",
			page: &grokipedia.Page{
				Title:   "NoDesc",
				Slug:    "nodesc",
				Content: "Just content.",
			},
			wantSub:    []string{"Title: NoDesc", "Slug: nodesc", "Just content."},
			notWantSub: []string{"Description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.FormatPage(tt.page)
			if err != nil {
				t.Fatalf("FormatPage() error = %v", err)
			}

			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("FormatPage() missing %q in:\n%s", sub, got)
				}
			}
			for _, sub := range tt.notWantSub {
				if strings.Contains(got, sub) {
					t.Errorf("FormatPage() should not contain %q in:\n%s", sub, got)
				}
			}
		})
	}
}

func TestTextFormatter_NoResults(t *testing.T) {
	f := NewText()
	got := f.NoResults()
	want := "No results found. Try different keywords.\n"
	if got != want {
		t.Errorf("NoResults() = %q, want %q", got, want)
	}
}

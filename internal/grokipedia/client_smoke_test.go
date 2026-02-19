package grokipedia

import (
	"os"
	"testing"
)

func TestClient_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping smoke tests in short mode")
	}

	if os.Getenv("GROKIR_E2E") != "1" {
		t.Skip("set GROKIR_E2E=1 to run smoke tests against real API")
	}

	client := NewClient()
	results, err := client.Search("kubernetes", 3, 0)
	if err != nil {
		t.Fatalf("Search() setup error = %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected search results, got empty")
	}

	top := results[0]

	t.Run("Search", func(t *testing.T) {
		if top.Title == "" {
			t.Error("expected non-empty title")
		}
		if top.Slug == "" {
			t.Error("expected non-empty slug")
		}
	})

	t.Run("GetPage", func(t *testing.T) {
		slug := top.Slug
		page, err := client.GetPage(slug, true, false)
		if err != nil {
			t.Fatalf("GetPage() error = %v", err)
		}
		if page == nil {
			t.Fatal("expected page, got nil")
		}
		if page.Title == "" {
			t.Error("expected non-empty title")
		}
		if page.Slug != slug {
			t.Errorf("slug = %q, want %q", page.Slug, slug)
		}
	})
}

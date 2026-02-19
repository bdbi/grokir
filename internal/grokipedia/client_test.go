package grokipedia

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestClient_Search(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		limit       int
		offset      int
		statusCode  int
		response    string
		wantLen     int
		wantErr     bool
		errContains string
	}{
		{
			name:       "success",
			query:      "kubernetes",
			limit:      5,
			offset:     0,
			statusCode: http.StatusOK,
			response:   `{"results":[{"slug":"k8s","title":"Kubernetes","snippet":"Container orchestrator","relevance_score":0.9,"view_count":"1000"}]}`,
			wantLen:    1,
			wantErr:    false,
		},
		{
			name:       "empty results",
			query:      "nonexistent",
			limit:      10,
			offset:     0,
			statusCode: http.StatusOK,
			response:   `{"results":[]}`,
			wantLen:    0,
			wantErr:    false,
		},
		{
			name:        "server error",
			query:       "test",
			limit:       10,
			offset:      0,
			statusCode:  http.StatusInternalServerError,
			response:    "internal server error",
			wantLen:     0,
			wantErr:     true,
			errContains: "500",
		},
		{
			name:        "invalid json",
			query:       "test",
			limit:       10,
			offset:      0,
			statusCode:  http.StatusOK,
			response:    `{"results":[{invalid}]}`,
			wantLen:     0,
			wantErr:     true,
			errContains: "decoding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/full-text-search" {
					t.Errorf("expected path /api/full-text-search, got %s", r.URL.Path)
				}
				if r.URL.Query().Get("query") != tt.query {
					t.Errorf("expected query %s, got %s", tt.query, r.URL.Query().Get("query"))
				}
				if tt.limit > 0 {
					if got := r.URL.Query().Get("limit"); got != strconv.Itoa(tt.limit) {
						t.Errorf("expected limit %d, got %s", tt.limit, got)
					}
				}
				if tt.offset > 0 {
					if got := r.URL.Query().Get("offset"); got != strconv.Itoa(tt.offset) {
						t.Errorf("expected offset %d, got %s", tt.offset, got)
					}
				}
				if ua := r.Header.Get("User-Agent"); ua != "grokir-test-agent" {
					t.Errorf("expected User-Agent grokir-test-agent, got %q", ua)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := &Client{
				BaseURL:   server.URL,
				UserAgent: "grokir-test-agent",
				HTTP:      server.Client(),
			}

			results, err := client.Search(tt.query, tt.limit, tt.offset)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(results) != tt.wantLen {
				t.Errorf("got %d results, want %d", len(results), tt.wantLen)
			}
		})
	}
}

func TestClient_GetPage(t *testing.T) {
	tests := []struct {
		name           string
		slug           string
		includeContent bool
		validateLinks  bool
		statusCode     int
		response       string
		wantTitle      string
		wantErr        bool
		errContains    string
	}{
		{
			name:           "success",
			slug:           "kubernetes",
			includeContent: true,
			validateLinks:  false,
			statusCode:     http.StatusOK,
			response:       `{"found":true,"page":{"title":"Kubernetes","slug":"kubernetes","description":"Container orchestrator","content":"Content here"}}`,
			wantTitle:      "Kubernetes",
			wantErr:        false,
		},
		{
			name:           "not found",
			slug:           "nonexistent",
			includeContent: true,
			validateLinks:  false,
			statusCode:     http.StatusNotFound,
			response:       "not found",
			wantTitle:      "",
			wantErr:        true,
			errContains:    "not found",
		},
		{
			name:           "found=false in response",
			slug:           "deleted",
			includeContent: true,
			validateLinks:  false,
			statusCode:     http.StatusOK,
			response:       `{"found":false,"page":null}`,
			wantTitle:      "",
			wantErr:        true,
			errContains:    "not found",
		},
		{
			name:           "invalid json",
			slug:           "test",
			includeContent: true,
			validateLinks:  false,
			statusCode:     http.StatusOK,
			response:       `{"found":true,"page":{invalid}}`,
			wantTitle:      "",
			wantErr:        true,
			errContains:    "decoding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/page" {
					t.Errorf("expected path /api/page, got %s", r.URL.Path)
				}
				if r.URL.Query().Get("slug") != tt.slug {
					t.Errorf("expected slug %s, got %s", tt.slug, r.URL.Query().Get("slug"))
				}
				if got := r.URL.Query().Get("includeContent"); got != strconv.FormatBool(tt.includeContent) {
					t.Errorf("expected includeContent %t, got %s", tt.includeContent, got)
				}
				if got := r.URL.Query().Get("validateLinks"); got != strconv.FormatBool(tt.validateLinks) {
					t.Errorf("expected validateLinks %t, got %s", tt.validateLinks, got)
				}
				if ua := r.Header.Get("User-Agent"); ua != "grokir-test-agent" {
					t.Errorf("expected User-Agent grokir-test-agent, got %q", ua)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := &Client{
				BaseURL:   server.URL,
				UserAgent: "grokir-test-agent",
				HTTP:      server.Client(),
			}

			page, err := client.GetPage(tt.slug, tt.includeContent, tt.validateLinks)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if page == nil {
				t.Fatal("expected page, got nil")
			}

			if page.Title != tt.wantTitle {
				t.Errorf("title = %q, want %q", page.Title, tt.wantTitle)
			}
		})
	}
}

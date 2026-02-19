package grokipedia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL   = "https://grokipedia.com"
	defaultUserAgent = "grokir/0.1 (Go Grokipedia CLI)"
)

// Client is a minimal HTTP client for the Grokipedia REST API.
type Client struct {
	BaseURL   string
	UserAgent string
	HTTP      *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:   defaultBaseURL,
		UserAgent: defaultUserAgent,
		HTTP:      http.DefaultClient,
	}
}

// SearchResult matches the structure returned by /api/full-text-search.
type SearchResult struct {
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Snippet         string   `json:"snippet"`
	RelevanceScore  float64  `json:"relevance_score"`
	ViewCount       int64    `json:"view_count,string"`
	TitleHighlights []string `json:"title_highlights"`
	SnippetHighlights []string `json:"snippet_highlights"`
}

type searchResponse struct {
	Results []SearchResult `json:"results"`
}

// Page represents a Grokipedia page from /api/page.
type Page struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type pageResponse struct {
	Found bool  `json:"found"`
	Page  *Page `json:"page"`
}

// Search performs a full-text search on Grokipedia.
func (c *Client) Search(query string, limit, offset int) ([]SearchResult, error) {
	if c.HTTP == nil {
		c.HTTP = http.DefaultClient
	}

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = "/api/full-text-search"

	q := u.Query()
	q.Set("query", query)
	if limit > 0 {
		q.Set("limit", fmt.Sprint(limit))
	}
	if offset > 0 {
		q.Set("offset", fmt.Sprint(offset))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed: %s", resp.Status)
	}

	var sr searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return sr.Results, nil
}

// GetPage retrieves a page by slug.
func (c *Client) GetPage(slug string, includeContent, validateLinks bool) (*Page, error) {
	if c.HTTP == nil {
		c.HTTP = http.DefaultClient
	}

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = "/api/page"

	q := u.Query()
	q.Set("slug", slug)
	q.Set("includeContent", fmt.Sprintf("%v", includeContent))
	q.Set("validateLinks", fmt.Sprintf("%v", validateLinks))
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("page not found: %s", slug)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get page failed: %s", resp.Status)
	}

	var pr pageResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	if !pr.Found || pr.Page == nil {
		return nil, fmt.Errorf("page not found: %s", slug)
	}

	return pr.Page, nil
}

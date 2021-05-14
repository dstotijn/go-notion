package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL       = "https://api.notion.com/v1"
	apiVersion    = "2021-05-13"
	clientVersion = "0.0.0"
)

// Client is used for HTTP requests to the Notion API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// ClientOption is used to override default client behavior.
type ClientOption func(*Client)

// NewClient returns a new Client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithHTTPClient overrides the default http.Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiKey))
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("User-Agent", "go-notion/"+clientVersion)

	if method == http.MethodPost || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// FindDatabaseByID fetches a database by ID.
// See: https://developers.notion.com/reference/get-database
func (c *Client) FindDatabaseByID(ctx context.Context, id string) (db Database, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/databases/"+id, nil)
	if err != nil {
		return Database{}, fmt.Errorf("notion: invalid URL: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Database{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Database{}, fmt.Errorf("notion: failed to find database: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&db)
	if err != nil {
		return Database{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return db, nil
}

// QueryDatabase returns database contents, with optional filters, sorts and pagination.
// See: https://developers.notion.com/reference/post-database-query
func (c *Client) QueryDatabase(ctx context.Context, id string, query DatabaseQuery) (result DatabaseQueryResponse, err error) {
	body := &bytes.Buffer{}

	err = json.NewEncoder(body).Encode(query)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to encode filter to JSON: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/databases/%v/query", id), body)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: invalid URL: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to find database: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return result, nil
}

// FindPageByID fetches a page by ID.
// See: https://developers.notion.com/reference/get-page
func (c *Client) FindPageByID(ctx context.Context, id string) (page Page, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/pages/"+id, nil)
	if err != nil {
		return Page{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("notion: failed to find page: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&page)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return page, nil
}

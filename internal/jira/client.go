package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client interface {
	SearchIssues(project, status string) ([]Issue, error)
	GetIssue(key string) (*Issue, error)
}

type HTTPClient struct {
	baseURL    string
	email      string
	apiToken   string
	httpClient *http.Client
}

func NewClient(baseURL, email, apiToken string) *HTTPClient {
	return &HTTPClient{
		baseURL:  baseURL,
		email:    email,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *HTTPClient) doRequest(method, endpoint string, query url.Values) ([]byte, error) {
	// Build URL by joining base URL and endpoint, handling trailing slashes
	baseURL := strings.TrimSuffix(c.baseURL, "/")
	fullURL := baseURL + endpoint
	if query != nil {
		fullURL += "?" + query.Encode()
	}

	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.email, c.apiToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *HTTPClient) SearchIssues(project, status string) ([]Issue, error) {
	// JQL: project keys work without quotes, status with spaces needs quotes
	jql := fmt.Sprintf(`project = %s AND status = "%s" ORDER BY updated DESC`, project, status)

	query := url.Values{}
	query.Set("jql", jql)
	query.Set("fields", "summary,status,issuetype,priority,assignee,reporter,created,updated")
	query.Set("maxResults", "50")

	body, err := c.doRequest(http.MethodGet, "/rest/api/3/search/jql", query)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Issues, nil
}

func (c *HTTPClient) GetIssue(key string) (*Issue, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s", key)

	query := url.Values{}
	query.Set("fields", "summary,status,issuetype,priority,assignee,reporter,created,updated,description")

	body, err := c.doRequest(http.MethodGet, endpoint, query)
	if err != nil {
		return nil, err
	}

	var issue Issue
	if err := json.Unmarshal(body, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &issue, nil
}

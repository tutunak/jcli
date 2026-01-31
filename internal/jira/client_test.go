package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPClient_SearchIssues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		// Verify auth header
		user, pass, ok := r.BasicAuth()
		if !ok || user != "test@example.com" || pass != "token123" {
			t.Error("invalid auth credentials")
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		jql := r.URL.Query().Get("jql")
		if jql == "" {
			t.Error("missing jql parameter")
		}

		result := SearchResult{
			Total: 2,
			Issues: []Issue{
				{Key: "TEST-1", Fields: IssueFields{Summary: "First issue", Status: Status{Name: "In Progress"}}},
				{Key: "TEST-2", Fields: IssueFields{Summary: "Second issue", Status: Status{Name: "In Progress"}}},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test@example.com", "token123")
	issues, err := client.SearchIssues("TEST", "In Progress")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(issues))
	}

	if issues[0].Key != "TEST-1" {
		t.Errorf("expected first issue key TEST-1, got %s", issues[0].Key)
	}
}

func TestHTTPClient_GetIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		issue := Issue{
			Key: "TEST-123",
			Fields: IssueFields{
				Summary:     "Test issue",
				Description: "This is a test",
				Status:      Status{Name: "In Progress"},
				IssueType:   Type{Name: "Task"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issue)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test@example.com", "token123")
	issue, err := client.GetIssue("TEST-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if issue.Key != "TEST-123" {
		t.Errorf("expected key TEST-123, got %s", issue.Key)
	}
	if issue.Fields.Summary != "Test issue" {
		t.Errorf("expected summary 'Test issue', got %q", issue.Fields.Summary)
	}
}

func TestHTTPClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"errorMessages":["Issue not found"]}`, http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test@example.com", "token123")
	_, err := client.GetIssue("NONEXISTENT-999")
	if err == nil {
		t.Error("expected error for non-existent issue")
	}
}

func TestMockClient(t *testing.T) {
	mock := NewMockClient()
	mock.AddIssue(Issue{
		Key: "MOCK-1",
		Fields: IssueFields{
			Summary: "Mock issue",
			Status:  Status{Name: "In Progress"},
		},
	})
	mock.AddIssue(Issue{
		Key: "MOCK-2",
		Fields: IssueFields{
			Summary: "Done issue",
			Status:  Status{Name: "Done"},
		},
	})

	t.Run("SearchIssues filters by status", func(t *testing.T) {
		issues, err := mock.SearchIssues("MOCK", "In Progress")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 1 {
			t.Errorf("expected 1 issue, got %d", len(issues))
		}
		if issues[0].Key != "MOCK-1" {
			t.Errorf("expected MOCK-1, got %s", issues[0].Key)
		}
	})

	t.Run("GetIssue returns issue by key", func(t *testing.T) {
		issue, err := mock.GetIssue("MOCK-1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if issue.Key != "MOCK-1" {
			t.Errorf("expected MOCK-1, got %s", issue.Key)
		}
	})

	t.Run("GetIssue returns error for unknown key", func(t *testing.T) {
		_, err := mock.GetIssue("UNKNOWN-999")
		if err == nil {
			t.Error("expected error for unknown key")
		}
		if _, ok := err.(*NotFoundError); !ok {
			t.Errorf("expected NotFoundError, got %T", err)
		}
	})
}

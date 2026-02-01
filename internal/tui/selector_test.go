package tui

import (
	"testing"

	"github.com/tutunak/jcli/internal/jira"
)

func TestNewSelector(t *testing.T) {
	s := NewSelector()
	if s == nil {
		t.Error("expected non-nil selector")
	}
}

func TestSelectIssue_EmptyList(t *testing.T) {
	s := NewSelector()
	_, err := s.SelectIssue([]jira.Issue{})
	if err == nil {
		t.Error("expected error for empty issue list")
	}
}

// Note: Interactive tests for SelectIssue and PromptCredentials
// would require mocking the terminal, which is complex.
// These are better tested through integration tests or manual testing.

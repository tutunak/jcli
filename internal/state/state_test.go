package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStateDir(t *testing.T) {
	t.Run("uses XDG_STATE_HOME if set", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "/custom/state")
		dir, err := StateDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "/custom/state/jcli"
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	})

	t.Run("falls back to ~/.local/state/jcli", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "")
		dir, err := StateDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, ".local", "state", "jcli")
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	})
}

func TestLoadAndSave(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_STATE_HOME", tmpDir)

	// Load should return empty state when file doesn't exist
	s, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading non-existent state: %v", err)
	}
	if s.HasCurrentIssue() {
		t.Error("expected no current issue in empty state")
	}

	// Set and save
	s.SetCurrentIssue("TEST-123", "Test summary")
	if err := s.Save(); err != nil {
		t.Fatalf("unexpected error saving state: %v", err)
	}

	// Reload and verify
	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading state: %v", err)
	}

	if !loaded.HasCurrentIssue() {
		t.Fatal("expected current issue to be set")
	}
	if loaded.CurrentIssue.Key != "TEST-123" {
		t.Errorf("expected key TEST-123, got %q", loaded.CurrentIssue.Key)
	}
	if loaded.CurrentIssue.Summary != "Test summary" {
		t.Errorf("expected summary 'Test summary', got %q", loaded.CurrentIssue.Summary)
	}
	if loaded.CurrentIssue.SelectedAt.IsZero() {
		t.Error("expected SelectedAt to be set")
	}
}

func TestSetCurrentIssue(t *testing.T) {
	s := &State{}

	s.SetCurrentIssue("PROJ-456", "My task")

	if s.CurrentIssue == nil {
		t.Fatal("expected CurrentIssue to be set")
	}
	if s.CurrentIssue.Key != "PROJ-456" {
		t.Errorf("expected key PROJ-456, got %q", s.CurrentIssue.Key)
	}
	if s.CurrentIssue.Summary != "My task" {
		t.Errorf("expected summary 'My task', got %q", s.CurrentIssue.Summary)
	}
}

func TestClearCurrentIssue(t *testing.T) {
	s := &State{}
	s.SetCurrentIssue("TEST-1", "Test")

	if !s.HasCurrentIssue() {
		t.Fatal("expected HasCurrentIssue to be true before clear")
	}

	s.ClearCurrentIssue()

	if s.HasCurrentIssue() {
		t.Error("expected HasCurrentIssue to be false after clear")
	}
}

func TestHasCurrentIssue(t *testing.T) {
	s := &State{}

	if s.HasCurrentIssue() {
		t.Error("expected HasCurrentIssue to be false for empty state")
	}

	s.CurrentIssue = &CurrentIssue{Key: "X-1"}
	if !s.HasCurrentIssue() {
		t.Error("expected HasCurrentIssue to be true when issue is set")
	}
}

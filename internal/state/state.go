package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CurrentIssue struct {
	Key        string    `json:"key"`
	Summary    string    `json:"summary"`
	SelectedAt time.Time `json:"selected_at"`
}

type State struct {
	CurrentIssue *CurrentIssue `json:"current_issue,omitempty"`
}

func StateDir() (string, error) {
	if xdgState := os.Getenv("XDG_STATE_HOME"); xdgState != "" {
		return filepath.Join(xdgState, "jcli"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".local", "state", "jcli"), nil
}

func StatePath() (string, error) {
	dir, err := StateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "state.json"), nil
}

func Load() (*State, error) {
	path, err := StatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &State{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

func (s *State) Save() error {
	path, err := StatePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

func (s *State) SetCurrentIssue(key, summary string) {
	s.CurrentIssue = &CurrentIssue{
		Key:        key,
		Summary:    summary,
		SelectedAt: time.Now(),
	}
}

func (s *State) ClearCurrentIssue() {
	s.CurrentIssue = nil
}

func (s *State) HasCurrentIssue() bool {
	return s.CurrentIssue != nil
}

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Defaults.Status != "In Progress" {
		t.Errorf("expected default status 'In Progress', got %q", cfg.Defaults.Status)
	}
}

func TestConfigDir(t *testing.T) {
	t.Run("uses XDG_CONFIG_HOME if set", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "/custom/config")
		dir, err := ConfigDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		expected := "/custom/config/jcli"
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	})

	t.Run("falls back to ~/.config/jcli", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		dir, err := ConfigDir()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		home, _ := os.UserHomeDir()
		expected := filepath.Join(home, ".config", "jcli")
		if dir != expected {
			t.Errorf("expected %q, got %q", expected, dir)
		}
	})
}

func TestLoadAndSave(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Load should return default config when file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading non-existent config: %v", err)
	}
	if cfg.Defaults.Status != "In Progress" {
		t.Errorf("expected default status, got %q", cfg.Defaults.Status)
	}

	// Save and reload
	cfg.Jira.URL = "https://test.atlassian.net"
	cfg.Jira.Email = "test@example.com"
	cfg.Jira.APIToken = "token123"
	cfg.Defaults.Project = "TEST"

	if err := cfg.Save(); err != nil {
		t.Fatalf("unexpected error saving config: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if loaded.Jira.URL != "https://test.atlassian.net" {
		t.Errorf("expected URL to persist, got %q", loaded.Jira.URL)
	}
	if loaded.Defaults.Project != "TEST" {
		t.Errorf("expected project to persist, got %q", loaded.Defaults.Project)
	}
}

func TestEnvOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create a config file
	cfg := DefaultConfig()
	cfg.Jira.URL = "https://file.atlassian.net"
	cfg.Jira.Email = "file@example.com"
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Set env vars
	t.Setenv("JIRA_URL", "https://env.atlassian.net")
	t.Setenv("JIRA_EMAIL", "env@example.com")
	t.Setenv("JIRA_API_TOKEN", "env-token")
	t.Setenv("JIRA_PROJECT", "ENVPROJ")
	t.Setenv("JIRA_STATUS", "Done")

	loaded, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loaded.Jira.URL != "https://env.atlassian.net" {
		t.Errorf("expected env URL override, got %q", loaded.Jira.URL)
	}
	if loaded.Jira.Email != "env@example.com" {
		t.Errorf("expected env email override, got %q", loaded.Jira.Email)
	}
	if loaded.Jira.APIToken != "env-token" {
		t.Errorf("expected env token override, got %q", loaded.Jira.APIToken)
	}
	if loaded.Defaults.Project != "ENVPROJ" {
		t.Errorf("expected env project override, got %q", loaded.Defaults.Project)
	}
	if loaded.Defaults.Status != "Done" {
		t.Errorf("expected env status override, got %q", loaded.Defaults.Status)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "empty config",
			cfg:     &Config{},
			wantErr: true,
		},
		{
			name: "missing email",
			cfg: &Config{
				Jira: JiraConfig{URL: "https://test.atlassian.net"},
			},
			wantErr: true,
		},
		{
			name: "missing token",
			cfg: &Config{
				Jira: JiraConfig{
					URL:   "https://test.atlassian.net",
					Email: "test@example.com",
				},
			},
			wantErr: true,
		},
		{
			name: "valid config",
			cfg: &Config{
				Jira: JiraConfig{
					URL:      "https://test.atlassian.net",
					Email:    "test@example.com",
					APIToken: "token",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHasProject(t *testing.T) {
	cfg := &Config{}
	if cfg.HasProject() {
		t.Error("expected HasProject() to return false for empty project")
	}

	cfg.Defaults.Project = "TEST"
	if !cfg.HasProject() {
		t.Error("expected HasProject() to return true for set project")
	}
}

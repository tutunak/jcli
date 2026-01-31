//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCLIIntegration(t *testing.T) {
	// Build the CLI
	tmpBin := filepath.Join(t.TempDir(), "jcli")
	cmd := exec.Command("go", "build", "-o", tmpBin, "../../.")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build: %v\n%s", err, output)
	}

	// Setup test directories
	configDir := t.TempDir()
	stateDir := t.TempDir()

	// Create mock Jira server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/rest/api/3/search"):
			json.NewEncoder(w).Encode(map[string]interface{}{
				"total": 2,
				"issues": []map[string]interface{}{
					{
						"key": "TEST-1",
						"fields": map[string]interface{}{
							"summary": "First test issue",
							"status":  map[string]string{"name": "In Progress"},
						},
					},
					{
						"key": "TEST-2",
						"fields": map[string]interface{}{
							"summary": "Second test issue",
							"status":  map[string]string{"name": "In Progress"},
						},
					},
				},
			})
		case strings.HasPrefix(r.URL.Path, "/rest/api/3/issue/"):
			key := strings.TrimPrefix(r.URL.Path, "/rest/api/3/issue/")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"key": key,
				"fields": map[string]interface{}{
					"summary": "Test issue " + key,
					"status":  map[string]string{"name": "In Progress"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Helper to run CLI
	runCLI := func(args ...string) (string, error) {
		cmd := exec.Command(tmpBin, args...)
		cmd.Env = append(os.Environ(),
			"XDG_CONFIG_HOME="+configDir,
			"XDG_STATE_HOME="+stateDir,
		)
		output, err := cmd.CombinedOutput()
		return string(output), err
	}

	// Test version
	t.Run("version", func(t *testing.T) {
		output, err := runCLI("version")
		if err != nil {
			t.Fatalf("version failed: %v", err)
		}
		if !strings.Contains(output, "jcli version") {
			t.Errorf("unexpected version output: %s", output)
		}
	})

	// Test help
	t.Run("help", func(t *testing.T) {
		output, err := runCLI("help")
		if err != nil {
			t.Fatalf("help failed: %v", err)
		}
		if !strings.Contains(output, "jcli - Jira CLI") {
			t.Errorf("unexpected help output: %s", output)
		}
	})

	// Test config project
	t.Run("config project", func(t *testing.T) {
		output, err := runCLI("config", "project", "TEST")
		if err != nil {
			t.Fatalf("config project failed: %v", err)
		}
		if !strings.Contains(output, "Default project set to: TEST") {
			t.Errorf("unexpected output: %s", output)
		}

		// Verify config file was created
		configFile := filepath.Join(configDir, "jcli", "config.yaml")
		data, err := os.ReadFile(configFile)
		if err != nil {
			t.Fatalf("failed to read config: %v", err)
		}

		var cfg map[string]interface{}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			t.Fatalf("failed to parse config: %v", err)
		}

		defaults := cfg["defaults"].(map[string]interface{})
		if defaults["project"] != "TEST" {
			t.Errorf("expected project TEST, got %v", defaults["project"])
		}
	})

	// Test config status
	t.Run("config status", func(t *testing.T) {
		output, err := runCLI("config", "status", "Done")
		if err != nil {
			t.Fatalf("config status failed: %v", err)
		}
		if !strings.Contains(output, "Default status filter set to: Done") {
			t.Errorf("unexpected output: %s", output)
		}
	})

	// Test config show
	t.Run("config show", func(t *testing.T) {
		output, err := runCLI("config", "show")
		if err != nil {
			t.Fatalf("config show failed: %v", err)
		}
		if !strings.Contains(output, "Project: TEST") {
			t.Errorf("unexpected output: %s", output)
		}
		if !strings.Contains(output, "Status: Done") {
			t.Errorf("unexpected output: %s", output)
		}
	})

	// Setup full config for task commands
	configFile := filepath.Join(configDir, "jcli", "config.yaml")
	cfg := map[string]interface{}{
		"jira": map[string]string{
			"url":       server.URL,
			"email":     "test@example.com",
			"api_token": "test-token",
		},
		"defaults": map[string]string{
			"project": "TEST",
			"status":  "In Progress",
		},
	}
	data, _ := yaml.Marshal(cfg)
	os.WriteFile(configFile, data, 0600)

	// Test task select with issue ID
	t.Run("task select by ID", func(t *testing.T) {
		output, err := runCLI("task", "select", "TEST-123")
		if err != nil {
			t.Fatalf("task select failed: %v\n%s", err, output)
		}
		if !strings.Contains(output, "Selected: TEST-123") {
			t.Errorf("unexpected output: %s", output)
		}
	})

	// Test task current
	t.Run("task current", func(t *testing.T) {
		output, err := runCLI("task", "current")
		if err != nil {
			t.Fatalf("task current failed: %v", err)
		}
		if !strings.Contains(output, "Current task: TEST-123") {
			t.Errorf("unexpected output: %s", output)
		}
	})

	// Test task branch
	t.Run("task branch", func(t *testing.T) {
		output, err := runCLI("task", "branch")
		if err != nil {
			t.Fatalf("task branch failed: %v", err)
		}
		output = strings.TrimSpace(output)
		if !strings.HasPrefix(output, "test-123-") {
			t.Errorf("branch should start with 'test-123-', got: %s", output)
		}
	})

	// Test task help
	t.Run("task help", func(t *testing.T) {
		output, err := runCLI("task", "help")
		if err != nil {
			t.Fatalf("task help failed: %v", err)
		}
		if !strings.Contains(output, "jcli task - Manage Jira tasks") {
			t.Errorf("unexpected output: %s", output)
		}
	})
}

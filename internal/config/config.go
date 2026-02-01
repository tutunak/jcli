package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type JiraConfig struct {
	URL      string `yaml:"url"`
	Email    string `yaml:"email"`
	APIToken string `yaml:"api_token"`
}

type Defaults struct {
	Project string `yaml:"project"`
	Status  string `yaml:"status"`
}

type Config struct {
	Jira     JiraConfig `yaml:"jira"`
	Defaults Defaults   `yaml:"defaults"`
}

func DefaultConfig() *Config {
	return &Config{
		Defaults: Defaults{
			Status: "In Progress",
		},
	}
}

func ConfigDir() (string, error) {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "jcli"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "jcli"), nil
}

func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg.applyEnvOverrides()
		return cfg, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.applyEnvOverrides()
	return cfg, nil
}

func (c *Config) applyEnvOverrides() {
	if url := os.Getenv("JIRA_URL"); url != "" {
		c.Jira.URL = url
	}
	if email := os.Getenv("JIRA_EMAIL"); email != "" {
		c.Jira.Email = email
	}
	if token := os.Getenv("JIRA_API_TOKEN"); token != "" {
		c.Jira.APIToken = token
	}
	if project := os.Getenv("JIRA_PROJECT"); project != "" {
		c.Defaults.Project = project
	}
	if status := os.Getenv("JIRA_STATUS"); status != "" {
		c.Defaults.Status = status
	}
}

func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) Validate() error {
	if c.Jira.URL == "" {
		return fmt.Errorf("jira.url is not configured")
	}
	if c.Jira.Email == "" {
		return fmt.Errorf("jira.email is not configured")
	}
	if c.Jira.APIToken == "" {
		return fmt.Errorf("jira.api_token is not configured (set via config or JIRA_API_TOKEN env var)")
	}
	return nil
}

func (c *Config) HasProject() bool {
	return c.Defaults.Project != ""
}

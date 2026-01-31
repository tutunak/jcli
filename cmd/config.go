package cmd

import (
	"fmt"
	"os"

	"github.com/dk/jcli/internal/config"
	"github.com/dk/jcli/internal/tui"
)

func executeConfig(args []string) error {
	if len(args) == 0 {
		printConfigUsage()
		return nil
	}

	switch args[0] {
	case "project":
		return executeConfigProject(args[1:])
	case "status":
		return executeConfigStatus(args[1:])
	case "credentials":
		return executeConfigCredentials(args[1:])
	case "show":
		return executeConfigShow()
	case "help", "--help", "-h":
		printConfigUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown config command: %s\n", args[0])
		printConfigUsage()
		return fmt.Errorf("unknown config command: %s", args[0])
	}
}

func printConfigUsage() {
	fmt.Println(`jcli config - Configure jcli settings

Usage:
  jcli config <command> [value]

Commands:
  project <key>     Set default Jira project
  status <name>     Set default status filter (default: "In Progress")
  credentials       Set Jira credentials interactively
  show              Show current configuration

Examples:
  jcli config project MYPROJ
  jcli config status "To Do"
  jcli config credentials
  jcli config show`)
}

func executeConfigProject(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("project key required")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Defaults.Project = args[0]

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Default project set to: %s\n", args[0])
	return nil
}

func executeConfigStatus(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("status name required")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cfg.Defaults.Status = args[0]

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Default status filter set to: %s\n", args[0])
	return nil
}

func executeConfigCredentials(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	selector := tui.NewSelector()
	url, email, token, err := selector.PromptCredentials()
	if err != nil {
		return err
	}

	cfg.Jira.URL = url
	cfg.Jira.Email = email
	cfg.Jira.APIToken = token

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("Credentials saved successfully.")
	fmt.Println("Note: API token is stored in the config file. You can also use JIRA_API_TOKEN environment variable.")
	return nil
}

func executeConfigShow() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	configPath, _ := config.ConfigPath()

	fmt.Println("Configuration:")
	fmt.Printf("  Config file: %s\n", configPath)
	fmt.Println()
	fmt.Println("Jira:")
	fmt.Printf("  URL: %s\n", maskEmpty(cfg.Jira.URL))
	fmt.Printf("  Email: %s\n", maskEmpty(cfg.Jira.Email))
	fmt.Printf("  API Token: %s\n", maskSecret(cfg.Jira.APIToken))
	fmt.Println()
	fmt.Println("Defaults:")
	fmt.Printf("  Project: %s\n", maskEmpty(cfg.Defaults.Project))
	fmt.Printf("  Status: %s\n", cfg.Defaults.Status)

	return nil
}

func maskEmpty(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}

func maskSecret(s string) string {
	if s == "" {
		return "(not set)"
	}
	if len(s) <= 4 {
		return "****"
	}
	return s[:4] + "****"
}

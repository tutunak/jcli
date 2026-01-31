package cmd

import (
	"fmt"
	"os"

	"github.com/dk/jcli/internal/config"
	"github.com/dk/jcli/internal/jira"
	"github.com/dk/jcli/internal/state"
	"github.com/dk/jcli/internal/tui"
)

func executeTaskSelect(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "Configuration error:", err)
		fmt.Fprintln(os.Stderr, "Run 'jcli config credentials' to set up your Jira credentials.")
		return err
	}

	if !cfg.HasProject() {
		fmt.Fprintln(os.Stderr, "Warning: No default project set. Run 'jcli config project <KEY>' to set one.")
		return fmt.Errorf("no project configured")
	}

	client := jira.NewClient(cfg.Jira.URL, cfg.Jira.Email, cfg.Jira.APIToken)
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// If issue ID provided, select it directly
	if len(args) > 0 {
		issueKey := args[0]
		return selectIssueByKey(client, st, issueKey)
	}

	// Interactive selection
	return selectIssueInteractive(client, st, cfg)
}

func selectIssueByKey(client jira.Client, st *state.State, issueKey string) error {
	issue, err := client.GetIssue(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get issue %s: %w", issueKey, err)
	}

	st.SetCurrentIssue(issue.Key, issue.Fields.Summary)
	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("Selected: %s - %s\n", issue.Key, issue.Fields.Summary)
	return nil
}

func selectIssueInteractive(client jira.Client, st *state.State, cfg *config.Config) error {
	issues, err := client.SearchIssues(cfg.Defaults.Project, cfg.Defaults.Status)
	if err != nil {
		return fmt.Errorf("failed to search issues: %w", err)
	}

	if len(issues) == 0 {
		fmt.Printf("No issues found in project %s with status %q\n", cfg.Defaults.Project, cfg.Defaults.Status)
		return nil
	}

	selector := tui.NewSelector()
	selected, err := selector.SelectIssue(issues)
	if err != nil {
		return err
	}

	st.SetCurrentIssue(selected.Key, selected.Fields.Summary)
	if err := st.Save(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	fmt.Printf("Selected: %s - %s\n", selected.Key, selected.Fields.Summary)
	return nil
}

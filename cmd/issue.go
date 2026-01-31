package cmd

import (
	"fmt"
	"os"
)

func executeIssue(args []string) error {
	if len(args) == 0 {
		printIssueUsage()
		return nil
	}

	switch args[0] {
	case "select":
		return executeIssueSelect(args[1:])
	case "current":
		return executeIssueCurrent(args[1:])
	case "branch":
		return executeIssueBranch(args[1:])
	case "help", "--help", "-h":
		printIssueUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown issue command: %s\n", args[0])
		printIssueUsage()
		return fmt.Errorf("unknown issue command: %s", args[0])
	}
}

func printIssueUsage() {
	fmt.Println(`jcli issue - Manage Jira issues

Usage:
  jcli issue <command> [flags]

Commands:
  select [issue-id]   Select an issue (interactive or by ID)
  current             Show current active issue
  branch              Generate branch name for current issue

Examples:
  jcli issue select              # Interactive selection from In Progress issues
  jcli issue select PROJ-123     # Select specific issue
  jcli issue current             # Show currently selected issue
  jcli issue branch              # Generate branch name for current issue`)
}

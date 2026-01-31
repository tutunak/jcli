package cmd

import (
	"fmt"
	"os"
)

func executeTask(args []string) error {
	if len(args) == 0 {
		printTaskUsage()
		return nil
	}

	switch args[0] {
	case "select":
		return executeTaskSelect(args[1:])
	case "current":
		return executeTaskCurrent(args[1:])
	case "branch":
		return executeTaskBranch(args[1:])
	case "help", "--help", "-h":
		printTaskUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown task command: %s\n", args[0])
		printTaskUsage()
		return fmt.Errorf("unknown task command: %s", args[0])
	}
}

func printTaskUsage() {
	fmt.Println(`jcli task - Manage Jira tasks

Usage:
  jcli task <command> [flags]

Commands:
  select [issue-id]   Select a task (interactive or by ID)
  current             Show current active task
  branch              Generate branch name for current task

Examples:
  jcli task select              # Interactive selection from In Progress tasks
  jcli task select PROJ-123     # Select specific issue
  jcli task current             # Show currently selected task
  jcli task branch              # Generate branch name for current task`)
}

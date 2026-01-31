package cmd

import (
	"fmt"
	"os"
)

var version = "dev"

func SetVersion(v string) {
	version = v
}

func Execute() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	switch os.Args[1] {
	case "version", "--version", "-v":
		fmt.Printf("jcli version %s\n", version)
		return nil
	case "help", "--help", "-h":
		printUsage()
		return nil
	case "task":
		return executeTask(os.Args[2:])
	case "config":
		return executeConfig(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func printUsage() {
	fmt.Println(`jcli - Jira CLI workflow management tool

Usage:
  jcli <command> [subcommand] [flags]

Commands:
  task      Manage Jira tasks
  config    Configure jcli settings
  version   Print version information
  help      Show this help message

Task Commands:
  jcli task select [issue-id]   Select a task (interactive or by ID)
  jcli task current             Show current active task
  jcli task branch              Generate branch name for current task

Config Commands:
  jcli config project <key>     Set default project
  jcli config status <name>     Set default status filter
  jcli config credentials       Set Jira credentials (interactive)

Use "jcli <command> --help" for more information about a command.`)
}

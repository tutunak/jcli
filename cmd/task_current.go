package cmd

import (
	"fmt"

	"github.com/dk/jcli/internal/state"
)

func executeTaskCurrent(args []string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if !st.HasCurrentIssue() {
		fmt.Println("No task currently selected.")
		fmt.Println("Use 'jcli task select' to select a task.")
		return nil
	}

	issue := st.CurrentIssue
	fmt.Printf("Current task: %s\n", issue.Key)
	fmt.Printf("Summary: %s\n", issue.Summary)
	fmt.Printf("Selected at: %s\n", issue.SelectedAt.Format("2006-01-02 15:04:05"))

	return nil
}

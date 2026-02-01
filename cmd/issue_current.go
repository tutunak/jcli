package cmd

import (
	"fmt"

	"github.com/tutunak/jcli/internal/state"
)

func executeIssueCurrent(args []string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if !st.HasCurrentIssue() {
		fmt.Println("No issue currently selected.")
		fmt.Println("Use 'jcli issue select' to select an issue.")
		return nil
	}

	issue := st.CurrentIssue
	fmt.Printf("Current issue: %s\n", issue.Key)
	fmt.Printf("Summary: %s\n", issue.Summary)
	fmt.Printf("Selected at: %s\n", issue.SelectedAt.Format("2006-01-02 15:04:05"))

	return nil
}

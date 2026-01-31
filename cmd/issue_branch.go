package cmd

import (
	"fmt"

	"github.com/dk/jcli/internal/branch"
	"github.com/dk/jcli/internal/state"
)

func executeIssueBranch(args []string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if !st.HasCurrentIssue() {
		fmt.Println("No issue currently selected.")
		fmt.Println("Use 'jcli issue select' to select an issue first.")
		return fmt.Errorf("no issue selected")
	}

	issue := st.CurrentIssue
	gen := branch.NewGenerator()
	branchName := gen.Generate(issue.Key, issue.Summary)

	fmt.Println(branchName)
	return nil
}

package cmd

import (
	"fmt"

	"github.com/dk/jcli/internal/branch"
	"github.com/dk/jcli/internal/state"
)

func executeTaskBranch(args []string) error {
	st, err := state.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	if !st.HasCurrentIssue() {
		fmt.Println("No task currently selected.")
		fmt.Println("Use 'jcli task select' to select a task first.")
		return fmt.Errorf("no task selected")
	}

	issue := st.CurrentIssue
	gen := branch.NewGenerator()
	branchName := gen.Generate(issue.Key, issue.Summary)

	fmt.Println(branchName)
	return nil
}

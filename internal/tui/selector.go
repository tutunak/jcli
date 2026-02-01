package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/tutunak/jcli/internal/jira"
)

type Selector struct{}

func NewSelector() *Selector {
	return &Selector{}
}

func (s *Selector) SelectIssue(issues []jira.Issue) (*jira.Issue, error) {
	if len(issues) == 0 {
		return nil, fmt.Errorf("no issues available to select")
	}

	options := make([]huh.Option[string], len(issues))
	issueMap := make(map[string]*jira.Issue)

	for i, issue := range issues {
		label := fmt.Sprintf("%s: %s", issue.Key, issue.Fields.Summary)
		if len(label) > 80 {
			label = label[:77] + "..."
		}
		options[i] = huh.NewOption(label, issue.Key)
		issueMap[issue.Key] = &issues[i]
	}

	var selectedKey string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select an issue").
				Options(options...).
				Value(&selectedKey),
		),
	)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("selection cancelled: %w", err)
	}

	selected, ok := issueMap[selectedKey]
	if !ok {
		return nil, fmt.Errorf("selected issue not found")
	}

	return selected, nil
}

func (s *Selector) PromptCredentials() (url, email, token string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Jira URL").
				Description("e.g., https://company.atlassian.net").
				Value(&url).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("URL is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("Email").
				Description("Your Jira account email").
				Value(&email).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("email is required")
					}
					return nil
				}),
			huh.NewInput().
				Title("API Token").
				Description("Create at https://id.atlassian.com/manage-profile/security/api-tokens").
				EchoMode(huh.EchoModePassword).
				Value(&token).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("API token is required")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return "", "", "", fmt.Errorf("credentials input cancelled: %w", err)
	}

	return url, email, token, nil
}

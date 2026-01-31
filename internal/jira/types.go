package jira

import "time"

type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string  `json:"summary"`
	Description string  `json:"description"`
	Status      Status  `json:"status"`
	IssueType   Type    `json:"issuetype"`
	Priority    *Priority `json:"priority,omitempty"`
	Assignee    *User   `json:"assignee,omitempty"`
	Reporter    *User   `json:"reporter,omitempty"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
}

type Status struct {
	Name string `json:"name"`
}

type Type struct {
	Name string `json:"name"`
}

type Priority struct {
	Name string `json:"name"`
}

type User struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

type SearchResult struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

type SelectedIssue struct {
	Key        string    `json:"key"`
	Summary    string    `json:"summary"`
	SelectedAt time.Time `json:"selected_at"`
}

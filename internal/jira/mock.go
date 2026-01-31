package jira

type MockClient struct {
	Issues     []Issue
	IssueByKey map[string]*Issue
	SearchErr  error
	GetErr     error
}

func NewMockClient() *MockClient {
	return &MockClient{
		IssueByKey: make(map[string]*Issue),
	}
}

func (m *MockClient) AddIssue(issue Issue) {
	m.Issues = append(m.Issues, issue)
	m.IssueByKey[issue.Key] = &issue
}

func (m *MockClient) SearchIssues(project, status string) ([]Issue, error) {
	if m.SearchErr != nil {
		return nil, m.SearchErr
	}

	var filtered []Issue
	for _, issue := range m.Issues {
		if issue.Fields.Status.Name == status {
			filtered = append(filtered, issue)
		}
	}
	return filtered, nil
}

func (m *MockClient) GetIssue(key string) (*Issue, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}

	issue, ok := m.IssueByKey[key]
	if !ok {
		return nil, &NotFoundError{Key: key}
	}
	return issue, nil
}

type NotFoundError struct {
	Key string
}

func (e *NotFoundError) Error() string {
	return "issue not found: " + e.Key
}

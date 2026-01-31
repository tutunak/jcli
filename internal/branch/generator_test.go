package branch

import (
	"strings"
	"testing"
)

func TestGenerator_Generate(t *testing.T) {
	// Use fixed random for predictable tests
	gen := NewGeneratorWithRand(func() int { return 847291 })

	tests := []struct {
		name     string
		issueKey string
		summary  string
		want     string
	}{
		{
			name:     "simple summary",
			issueKey: "PROJ-123",
			summary:  "Add user authentication",
			want:     "PROJ-123-add-user-authentication-847291",
		},
		{
			name:     "summary with special characters",
			issueKey: "TEST-456",
			summary:  "Fix bug: user can't login (urgent!)",
			want:     "TEST-456-fix-bug-user-can-t-login-urgent-847291",
		},
		{
			name:     "summary with numbers",
			issueKey: "DEV-789",
			summary:  "Upgrade to v2.0",
			want:     "DEV-789-upgrade-to-v2-0-847291",
		},
		{
			name:     "summary with extra spaces",
			issueKey: "TASK-1",
			summary:  "  Multiple   spaces   here  ",
			want:     "TASK-1-multiple-spaces-here-847291",
		},
		{
			name:     "empty summary",
			issueKey: "PROJ-999",
			summary:  "",
			want:     "PROJ-999--847291",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gen.Generate(tt.issueKey, tt.summary)
			if got != tt.want {
				t.Errorf("Generate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGenerator_GenerateWithRandomNumber(t *testing.T) {
	gen := NewGenerator()
	branch := gen.Generate("TEST-1", "Test")

	// Should start with issue key preserving original case
	if !strings.HasPrefix(branch, "TEST-1-") {
		t.Errorf("branch should start with 'TEST-1-', got %q", branch)
	}

	// Should end with a number
	parts := strings.Split(branch, "-")
	lastPart := parts[len(parts)-1]
	for _, c := range lastPart {
		if c < '0' || c > '9' {
			t.Errorf("last part should be numeric, got %q", lastPart)
			break
		}
	}
}

func TestNormalizeSummary(t *testing.T) {
	tests := []struct {
		name    string
		summary string
		want    string
	}{
		{
			name:    "lowercase conversion",
			summary: "UPPERCASE Text",
			want:    "uppercase-text",
		},
		{
			name:    "special characters replaced",
			summary: "hello@world#test$123",
			want:    "hello-world-test-123",
		},
		{
			name:    "multiple hyphens collapsed",
			summary: "hello---world",
			want:    "hello-world",
		},
		{
			name:    "leading and trailing hyphens trimmed",
			summary: "---hello---",
			want:    "hello",
		},
		{
			name:    "unicode characters",
			summary: "café résumé naïve",
			want:    "cafe-resume-naive",
		},
		{
			name:    "long summary truncated",
			summary: "This is a very long summary that should be truncated to avoid creating branch names that are too long for most git systems",
			want:    "this-is-a-very-long-summary-that-should-be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeSummary(tt.summary)
			if got != tt.want {
				t.Errorf("normalizeSummary(%q) = %q, want %q", tt.summary, got, tt.want)
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{1, "1"},
		{123, "123"},
		{999999, "999999"},
		{1000000, "0"}, // wraps around
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatNumber(tt.n)
			if got != tt.want {
				t.Errorf("formatNumber(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

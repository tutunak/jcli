# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
make build          # Build binary with version from git tags
make test           # Run tests with race detection and coverage
make lint           # Run golangci-lint
make test-integration  # Run integration tests (requires Jira credentials)
go test ./internal/branch/...  # Run tests for a specific package
```

## Architecture

jcli is a CLI tool for Jira workflow management. It uses manual argument parsing (no CLI framework).

### Package Structure

- **cmd/** - Command handlers. `root.go` routes to subcommands via switch statement on `os.Args[1]`. Each command file (e.g., `issue_select.go`) contains one handler function.

- **internal/jira/** - Jira API client using REST API v3. `Client` interface allows mocking. `HTTPClient` implements actual API calls with Basic Auth. Note: descriptions use `json.RawMessage` because Jira v3 returns ADF format, not strings.

- **internal/config/** - YAML config at `~/.config/jcli/config.yaml`. Supports `JIRA_API_TOKEN` env var override.

- **internal/state/** - JSON state at `~/.local/state/jcli/state.json`. Tracks currently selected issue.

- **internal/branch/** - Branch name generator. Format: `<ISSUE-KEY>-<normalized-summary>-<random>`. Issue key preserves original case; summary is lowercased.

- **internal/tui/** - Interactive issue selector using `github.com/charmbracelet/huh`.

### Key Design Decisions

- No CLI framework (cobra, urfave/cli) - uses manual `os.Args` parsing
- `Client` interface in jira package enables `MockClient` for testing
- XDG Base Directory paths for config and state
- JQL query includes `assignee = currentUser()` to show only user's assigned issues

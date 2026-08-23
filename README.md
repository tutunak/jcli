# jcli

A command-line tool for Jira workflow management. Quickly select issues, track your current work, and generate consistent branch names.

**Website:** <https://jcli.tutunak.com/>

## Features

- **Interactive issue selection** - Browse and select from your assigned "In Progress" issues
- **Direct issue selection** - Select any issue by its key
- **Current issue tracking** - Keep track of what you're working on
- **Branch name generation** - Generate consistent, readable branch names from issue keys and summaries
- **XDG-compliant configuration** - Config stored in `~/.config/jcli/`

## Installation

### From source

```bash
git clone https://github.com/tutunak/jcli.git
cd jcli
go build -o jcli .
sudo mv jcli /usr/local/bin/
```

### Using Go

```bash
go install github.com/tutunak/jcli@latest
```

## Configuration

jcli uses a YAML configuration file located at `~/.config/jcli/config.yaml`.

### Initial Setup

1. **Set your Jira credentials:**

```bash
jcli config credentials
```

This will interactively prompt for:
- Jira URL (e.g., `https://yourcompany.atlassian.net`)
- Email address
- API token

2. **Set your default project:**

```bash
jcli config project YOUR_PROJECT_KEY
```

3. **Optionally, change the default status filter:**

```bash
jcli config status "In Progress"
```

### Getting a Jira API Token

1. Go to <https://id.atlassian.com/manage-profile/security/api-tokens>
2. Click "Create API token"
3. Give it a name (e.g., "jcli")
4. Copy the generated token

### Configuration File Format

```yaml
jira:
  url: https://yourcompany.atlassian.net
  email: your.email@company.com
  api_token: your_api_token_here

defaults:
  project: PROJ
  status: In Progress
```

### Environment Variables

Any config value can be overridden by an environment variable:

| Variable         | Overrides          |
|------------------|--------------------|
| `JIRA_URL`       | `jira.url`         |
| `JIRA_EMAIL`     | `jira.email`       |
| `JIRA_API_TOKEN` | `jira.api_token`   |
| `JIRA_PROJECT`   | `defaults.project` |
| `JIRA_STATUS`    | `defaults.status`  |

```bash
export JIRA_API_TOKEN=your_api_token_here
```

`XDG_CONFIG_HOME` and `XDG_STATE_HOME` relocate the config and state directories.

## Usage

### Select an Issue

**Interactive selection** - Shows your assigned issues with "In Progress" status:

```bash
jcli issue select
```

Use arrow keys to navigate and Enter to select.

**Direct selection** - Select a specific issue by key:

```bash
jcli issue select PROJ-123
```

### View Current Issue

Display the currently selected issue:

```bash
jcli issue current
```

Output:
```text
Current issue: PROJ-123
Summary: Implement user authentication
Selected at: 2024-01-15 10:30:00
```

### Generate Branch Name

Generate a branch name for the current issue:

```bash
jcli issue branch
```

Output:
```
PROJ-123-implement-user-authentication-847291
```

The branch name format is: `<ISSUE-KEY>-<normalized-summary>-<random-number>`

- Issue key is preserved in original case (uppercase)
- Summary is converted to lowercase
- Special characters are replaced with hyphens
- Multiple hyphens are collapsed
- Long summaries are truncated at 50 characters
- Random number (0-999999) ensures uniqueness

### Create a Git Branch

Combine with git to create and checkout a new branch:

```bash
git checkout -b $(jcli issue branch)
```

## Shell Aliases

Add these to `~/.bashrc` or `~/.zshrc` to branch off the selected issue in one word.
They work unchanged in bash and zsh.

```bash
# Branch off the issue currently selected in jcli.
gcj() {
  local branch
  branch=$(jcli issue branch) || {
    printf '%s\n' "$branch" >&2
    return 1
  }
  git checkout -b "$branch"
}

# Go back to the branch you already made for the current issue,
# creating one only if it does not exist yet.
gsj() {
  local line key existing
  line=$(jcli issue current | head -n 1)
  case $line in
    "Current issue: "*) key=${line#Current issue: } ;;
    *) printf '%s\n' "$line" >&2; return 1 ;;
  esac
  existing=$(git branch --list "${key}-*" \
    --format='%(refname:short)' | head -n 1)
  if [ -n "$existing" ]; then
    git checkout "$existing"
  else
    gcj
  fi
}

# Shorthands for the two commands you type most.
alias jis='jcli issue select'
alias jic='jcli issue current'
```

Reload with `source ~/.zshrc`, then:

```bash
jis          # pick an issue
gcj          # create and switch to a branch for it
gsj          # later, hop back to that branch
```

**Why functions and not a plain alias.** With no issue selected, `jcli issue branch`
prints its hint on *stdout* and exits 1. Command substitution keeps the stdout but
discards the exit status, so

```bash
alias gcj='git checkout -b $(jcli issue branch)'   # don't
```

hands git the hint text and fails with `fatal: Cannot update paths and switch to
branch 'No' at the same time.` The functions above check the exit status first and
forward the real message to stderr.

**Why `gsj` exists.** `jcli issue branch` appends a fresh random suffix on every
call, so running `gcj` twice for the same issue gives you two branches. `gsj` looks
for a branch already matching `<ISSUE-KEY>-*` and switches to it instead.

## Commands Reference

### Root Commands

| Command        | Description               |
|----------------|---------------------------|
| `jcli help`    | Show help message         |
| `jcli version` | Print version information |

### Issue Commands

| Command                   | Description                                              |
|---------------------------|----------------------------------------------------------|
| `jcli issue select`       | Interactive selection from assigned "In Progress" issues |
| `jcli issue select <KEY>` | Select a specific issue by key                           |
| `jcli issue current`      | Show currently selected issue                            |
| `jcli issue branch`       | Generate branch name for current issue                   |

### Config Commands

| Command                     | Description                        |
|-----------------------------|------------------------------------|
| `jcli config credentials`   | Set Jira credentials interactively |
| `jcli config project <KEY>` | Set default project key            |
| `jcli config status <NAME>` | Set default status filter          |

## Workflow Example

```bash
# One-time setup
jcli config credentials
jcli config project MYPROJ

# Daily workflow
jcli issue select                    # Pick an issue to work on
jcli issue current                   # Verify selection
git checkout -b $(jcli issue branch) # Create feature branch

# Or select and branch in one go
jcli issue select MYPROJ-456
git checkout -b $(jcli issue branch)
```

## File Locations

| File   | Location                         | Purpose                       |
|--------|----------------------------------|-------------------------------|
| Config | `~/.config/jcli/config.yaml`     | Jira credentials and defaults |
| State  | `~/.local/state/jcli/state.json` | Current issue tracking        |

## Development

### Prerequisites

- Go 1.21 or later

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

## License

MIT

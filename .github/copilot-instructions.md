# Tinycare-tui

Tinycare-tui is a Go terminal UI application that displays git commits from the last 24 hours and week, current weather, self-care advice, and todo list tasks. It's built with Go 1.20+ using the tview library for the terminal interface.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Bootstrap, Build, and Test the Repository
Run these commands in sequence to set up and validate the codebase:

```bash
# Install dependencies (takes ~8-10 seconds)
go mod download

# Build all packages (takes ~15-20 seconds)
go build ./...

# Install the binary (takes ~1 second)
go install ./cmd/tinycare-tui

# Run tests with coverage (takes ~4-7 seconds) - NEVER CANCEL
go test -v -cover ./internal/... -coverprofile coverage.out -coverpkg ./internal/...
```

**CRITICAL TIMING NOTES:**
- Build takes 15-20 seconds. NEVER CANCEL. Set timeout to 60+ seconds minimum.
- Tests take 4-7 seconds. NEVER CANCEL. Set timeout to 30+ seconds minimum.
- One test may fail (jokes API test) - this is expected and not a blocker.

### Code Quality and Linting
Always run these commands before submitting changes:

```bash
# Format code (takes <1 second)
go fmt ./...

# Vet code for issues (takes ~3 seconds)
go vet ./...
```

### Run the Application

#### Basic Run (Minimal Configuration)
```bash
# Create a test todo file for demonstration
echo -e "- [ ] Test task 1\n- [x] Completed task\n- [ ] Another task" > /tmp/test_todos.txt

# Run with minimal environment variables
TINYCARE_POSTAL_CODE=10001 TODO_FILE=/tmp/test_todos.txt tinycare-tui
```

#### Full Configuration Run
```bash
# Set all possible environment variables for complete functionality
export GITHUB_TOKEN=your_github_token_here
export OPEN_WEATHER_MAP_API_KEY=your_openweather_api_key_here
export TODOIST_TOKEN=your_todoist_token_here
export TINYCARE_POSTAL_CODE=10001
export TINYCARE_WORKSPACE=/path/to/git/repos,/another/path/to/repos
export TODO_FILE=/path/to/todolist/file

# Run the application
tinycare-tui
```

#### Test with Local Git Repository
```bash
# Create a test git repository with commits
cd /tmp
git init test_repo
cd test_repo
git config user.name "Test User"
git config user.email "test@test.com"
echo "# Test repo" > README.md
git add . && git commit -m "Initial commit"
echo "# Updated" >> README.md
git add . && git commit -m "Second commit"

# Run tinycare-tui with the test repository
TINYCARE_POSTAL_CODE=10001 TODO_FILE=/tmp/test_todos.txt TINYCARE_WORKSPACE=/tmp/test_repo tinycare-tui
```

## Validation

### Manual Validation Scenarios
After making changes, ALWAYS test these scenarios:

1. **Basic Startup Test:**
   - Run with minimal environment variables
   - Verify the application starts and shows the TUI interface
   - Check that appropriate error messages appear for missing API keys

2. **Git Integration Test:**
   - Create a test git repository with recent commits
   - Set TINYCARE_WORKSPACE to the test repository path
   - Verify commits appear in both Daily and Weekly commit panels

3. **Todo File Test:**
   - Create a test todo file with mixed completed/incomplete tasks
   - Set TODO_FILE environment variable
   - Verify tasks appear correctly in the Today's Tasks panel

4. **UI Navigation Test:**
   - Start the application
   - Press 'r' to refresh - verify self-care message changes
   - Press 'q' to quit - verify application exits cleanly
   - Press Tab to navigate between panels (if focus highlighting works)

### Expected Application Behavior
- **Startup:** Application loads within 2-3 seconds
- **Display:** Shows 4 main panels (Daily Commits, Weekly Commits, Weather, Self Care, Joke, Today's Tasks)
- **Controls:** 
  - 'q' exits the application
  - 'r' refreshes all content
  - Tab navigates between focusable panels
- **Auto-refresh:** Content refreshes every 30 seconds automatically
- **Error Handling:** Missing environment variables show helpful error messages

### CI Requirements
The GitHub Actions CI (.github/workflows/coverage.yml) will:
- Build the application with `go install ./...`
- Run tests with coverage requirement
- Generate coverage badge (currently 82.3%)

Always ensure your changes maintain or improve test coverage.

## Common Tasks

### Repository Structure
```
tinycare-tui/
├── cmd/tinycare-tui/        # Main application entry point
│   └── main.go             # Application setup and UI initialization
├── internal/               # Internal packages
│   ├── apis/              # External API integrations (GitHub, weather, jokes, Todoist)
│   ├── local/             # Local file and git repository handling
│   ├── ui/                # Terminal UI components and views
│   └── utils/             # Utility functions (formatting, etc.)
├── .github/workflows/     # CI/CD configuration
├── go.mod                 # Go module definition
└── go.sum                 # Go module checksums
```

### Key Files to Monitor
- `cmd/tinycare-tui/main.go` - Main application logic and UI setup
- `internal/apis/` - When working with external integrations
- `internal/local/git.go` - When working with git repository scanning
- `internal/ui/views.go` - When working with UI components
- `.github/workflows/coverage.yml` - CI configuration

### Environment Variables Reference
```bash
# Required for basic functionality
TINYCARE_POSTAL_CODE=12345        # For weather display

# Git commit sources (choose one)
TINYCARE_WORKSPACE=/path/to/repos  # Local git repositories
GITHUB_TOKEN=abc123               # GitHub API access

# Task sources (choose one)  
TODO_FILE=/path/to/file           # Local todo file
TODOIST_TOKEN=abc123              # Todoist API access

# Optional
OPEN_WEATHER_MAP_API_KEY=abc123   # Weather information
```

### Common Commands Output
```bash
# Repository root listing
$ ls -la
total 208
drwxr-xr-x 6 runner docker   4096 .
drwxr-xr-x 3 runner docker   4096 ..
drwxr-xr-x 7 runner docker   4096 .git
drwxr-xr-x 3 runner docker   4096 .github
-rw-r--r-- 1 runner docker    436 .gitignore
-rw-r--r-- 1 runner docker   1071 LICENSE
-rw-r--r-- 1 runner docker   3858 README.md
drwxr-xr-x 3 runner docker   4096 cmd
-rw-r--r-- 1 runner docker   1794 go.mod
-rw-r--r-- 1 runner docker  14719 go.sum
-rw-r--r-- 1 runner docker 152555 image.png
drwxr-xr-x 6 runner docker   4096 internal

# Go module info
$ cat go.mod
module github.com/DMcP89/tinycare-tui

go 1.20

require (
    github.com/PuerkitoBio/goquery v1.8.1
    github.com/gdamore/tcell/v2 v2.6.0
    github.com/go-git/go-git/v5 v5.7.0
    github.com/rivo/tview v0.0.0-20230530133550-8bd761dda819
)
```

### Troubleshooting
- **"object not found" error:** The git repository path in TINYCARE_WORKSPACE doesn't exist or isn't a valid git repo
- **Missing environment variable messages:** Expected behavior when API keys aren't provided
- **Test failures in jokes_test.go:** Known issue with external API dependency, not a blocker
- **Build issues:** Ensure Go 1.20+ is installed and go mod download has been run

### Development Tips
- Always test with both local git repositories and external APIs when possible
- The application is designed to handle missing environment variables gracefully
- Self-care messages and jokes are randomized on each refresh
- Use the test scenarios above to validate any UI or functional changes
- Coverage should remain above 80% for the awesome-go listing requirement
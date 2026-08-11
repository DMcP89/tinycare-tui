# AGENTS.md - Core Configuration for Tinycare-tui

## Project Overview
Tinycare-tui is a Go-based TUI application that provides:
- Git commit summaries (daily/weekly) from local or remote repos.
- Current weather information.
- Self-care affirmations and jokes.
- A Todoist or local text file task list.

## Key Technical Context
- **Frameworks**: Uses `tview` for the UI layout and `tcell` for terminal handling.
- **Concurrency**: Multiple "views" (Weather, Jokes, etc.) are updated in concurrent goroutines to prevent UI blocking during API calls or local disk I/O.
- **Data Sources**:
  - Internal & Local: Handled in `internal/local`.
  - External APIs: Handled in `internal/apis` using `internal/utils` for common request logic.

## Development workflow
- **Running**: `go run ./cmd/tinycare-tui` or simply `tinycare-tui` if installed.
- **Development requirement**: Environment variables are crucial. A missing variable will trigger a fallback message in the UI:
  - `GITHUB_TOKEN`: For GitHub data.
  - `TINYCARE_POSTAL_CODE`: For Weather.
  - `TODOIST_TOKEN` or `TODO_FILE`: For tasks.
  - `TASKWARRIOR`: A non-empty value (e.g. `1`) to pull tasks from a local Taskwarrior install via the `task` CLI.
  - `TINYCARE_WORKSPACE`: A comma-separated list of local paths to scan for Git repos.
- **Testing**: Run tests using standard Go tools (e.g., `go test ./...`). Note that some mock data exists in `internal/apis/testdata`.

## Rules & Conventions
- **View Updates**: Use the `GetTextForView` helper in `cmd/tinycare-tui/main.go` to handle optional environment variables gracefully before updating UI elements.
- **State Management**: The application uses a loop with `time.Sleep` based on `TINYCARE_REFRESH_INTERVAL`.

## Critical paths
- `internal/local/git.go`: Core logic for parsing local repository commit history.
- `internal/apis/*.go`: Integration points for external services.
- `cmd/tinycare-tui/main.go`: Main orchestration of the UI layout and refresh loops.

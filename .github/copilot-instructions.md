# Copilot Instructions for tscgit

## Project Overview
`tscgit` is an interactive command-line Git coach that teaches Git workflows through two modes:
- **Verification lessons** (`verify` command): Interactive checks against existing Git repositories using Bubble Tea UI
- **Run scripts** (`run` command): Scripted exercises that execute and validate shell commands step-by-step

## Architecture Pattern
This project follows a **registry-based plugin architecture**:

```
cmd/tscgit/main.go           # CLI entry point with command routing
internal/lessons/            # Verification lesson registry & gitutil helpers
internal/run/                # Run script registry & command executor  
internal/ui/{verify,run}/    # Bubble Tea UI models for each mode
internal/gitutil/            # Git repository wrapper with common operations
```

**Key principle**: Lessons and scripts self-register in `init()` functions, making the system easily extensible.

## Adding New Verification Lessons
Lessons live in `internal/lessons/`. Follow this pattern in `defaults.go` or new files:

```go
func init() {
    Must(Register(&Lesson{
        ID: "my-lesson",
        Title: "Short descriptive title",  
        Description: "Longer explanation for students",
        Checks: []Check{
            {
                ID: "check-id", 
                Title: "User-facing check name",
                Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
                    // Use repo.HasBranch(), repo.CommitCount(), repo.FileExists(), etc.
                    return CheckResult{Passed: true, Message: "Success message"}
                },
            },
        },
    }))
}
```

**gitutil.Repository** provides Git operations: `CurrentBranch()`, `HasBranch()`, `CommitCount()`, `LastCommitMessage()`, `CommitsAhead()`, `FileExists()`, `HasRemote()`. All methods take `context.Context`.

## Adding New Run Scripts  
Scripts live in `internal/run/scripts.go`. Register with expected commands and outputs:

```go
Register(&Script{
    ID: "42",
    Title: "Check something", 
    Description: "What this script validates",
    Steps: []Step{
        {
            Command: "git status",
            ExpectExitCode: 0,
            ExpectStdout: []string{"clean"}, // All strings must be present in output
        },
    },
})
```

**Shell behavior**: Windows uses PowerShell (`pwsh` or `powershell`), Unix uses `sh -c`. Set `ExpectExitCode: -1` to skip exit code validation. All strings in `ExpectStdout` must be present in command output.

## UI Architecture 
Both verification and run modes use **Bubble Tea** with similar patterns:
- Models in `internal/ui/{verify,run}/model.go`
- Async execution with progress spinners
- Results streamed as Tea messages (`checkResultMsg`, etc.)
- Exit codes: 0 = success, 1 = error, 2 = checks failed

## Development Workflow
```powershell
# Standard Go development (requires Go 1.25+)
gofmt -w .
go test ./...
go install ./cmd/tscgit

# Test the CLI
tscgit lessons                    # List available lessons/scripts
tscgit verify init-basics         # Run verification lesson
tscgit run 0                      # Run script lesson
```

## Testing Patterns
- Unit tests for core logic (executor, lessons registry)
- Use `gitutil.Repository` mock-friendly design for verification tests
- Shell command tests use real commands with expected outputs
- No integration tests - relies on manual CLI testing

## Module Dependencies
- **Bubble Tea ecosystem**: UI framework (charmbracelet/bubbletea, lipgloss, bubbles)
- **Standard library**: Core Git operations via `os/exec` 
- **No external Git libraries**: All Git operations through CLI commands for maximum compatibility

## Key Conventions
- Error handling: Return `CheckResult{Err: err}` for verification errors vs. `CheckResult{Passed: false}` for user issues
- Context usage: All Git operations accept `context.Context` for timeout/cancellation 
- Registration pattern: Use `Must()` wrapper for `init()` registrations to panic on conflicts
- Cross-platform: Shell detection in `defaultShell()`, path handling with `filepath.Join()`
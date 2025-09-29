# tscgit

`tscgit` is an interactive command-line coach for practicing Git workflows. It ships with ready-to-use lessons and a Bubble Tea-powered verification UI that guides students through checks in style.

## Installation

```powershell
# From the repository root
go install ./cmd/tscgit
```

This installs the `tscgit` binary into your Go bin directory (for Windows PowerShell users, typically `%USERPROFILE%\go\bin`). Make sure that directory is on your `PATH`.

## Usage

List available lessons:

```powershell
tscgit lessons
```

Verify your progress against a lesson:

```powershell
tscgit verify init-basics
```

Use `-path` to point at a different repository:

```powershell
tscgit verify -path C:\path\to\practice-repo init-basics
```

The verifier launches a terminal UI that streams each check with spinners, colorful feedback, and a summary once complete. Exit with `q` or `Ctrl+C` at any time.

Run scripted environment checks (the "Git Started" track):

```powershell
tscgit run 0
```

Each run script executes one or more commands (via PowerShell on Windows or sh on Unix) and validates exit codes and stdout to mimic an instructor walking through setup tasks.

### Git Started run scripts

| ID  | Title                        | What it checks |
|-----|------------------------------|----------------|
| `0` | Test your CLI                | Ensures `tscgit` can launch shell commands and echoes a sample phrase. |
| `1` | Install Git                  | Confirms `git --version` is available. |
| `2` | Configure Git identity       | Verifies `git config` global `user.name`, `user.email`, and default branch (`master`). |
| `3` | Initialize a repository      | Checks you're inside the `webflyx` repo and `.git` files exist. |
| `4a`| Track contents.md            | Validates the file exists with `# contents` prior to staging. |
| `4b`| Stage contents.md            | Confirms `contents.md` is staged. |
| `5` | Commit contents.md           | Looks for commit `A: add contents.md`. |
| `6a`| Inspect commit metadata      | Reads `catfileout.txt` for commit metadata. |
| `6b`| Inspect blob contents        | Ensures `blobfile.txt` mirrors blob output. |
| `7` | Create titles.md             | Verifies the `B:` commit and `titles.md` contents. |
| `8` | Switch default branch        | Checks global default branch and locals are `main`. |
| `9` | Create add_classics branch   | Confirms branch creation and checkout. |
| `10`| Commit classics.csv          | Looks for the `C:` commit while on `add_classics`. |
| `11a`| Prepare merge history       | Inspects the pre-merge graph for commits `A:` through `D:`. |
| `11b`| Merge add_classics          | Confirms merge commit `E:` with multiple parents. |

## Adding new verification lessons

Lessons live in `internal/lessons`. Each lesson bundles multiple checks:

```go
lessons.Must(lessons.Register(&lessons.Lesson{
    ID:          "my-lesson",
    Title:       "Awesome Git Flow",
    Description: "Introduce your students to rebase workflows.",
    Checks: []lessons.Check{
        {
            ID:    "has-develop",
            Title: "Create develop branch",
            Verify: func(ctx context.Context, repo *gitutil.Repository) lessons.CheckResult {
                ok, err := repo.HasBranch(ctx, "develop")
                if err != nil {
                    return lessons.CheckResult{Err: err}
                }
                if !ok {
                    return lessons.CheckResult{Passed: false, Message: "Branch 'develop' not found."}
                }
                return lessons.CheckResult{Passed: true, Message: "develop branch looks good."}
            },
        },
    },
}))
```

Checks receive a cancellable context plus a `gitutil.Repository` helper that wraps common Git queries.

## Adding new run scripts

Run scripts live in `internal/run`. Each script lists the exact commands students should have run and the expected outputs:

```go
runlesson.Register(&runlesson.Script{
    ID:          "12",
    Title:       "Check remotes",
    Description: "Ensure origin exists and points to GitHub.",
    Steps: []runlesson.Step{
        {
            Command:        "git remote",
            ExpectExitCode: 0,
            ExpectStdout:   []string{"origin"},
        },
        {
            Command:        "git remote get-url origin",
            ExpectExitCode: 0,
            ExpectStdout:   []string{"github.com"},
        },
    },
})
```

Commands are executed inside the user's shell (PowerShell on Windows, `sh` elsewhere). The runner enforces exit codes (set `ExpectExitCode` to `-1` to skip) and validates that every string listed in `ExpectStdout` is present in command output.

## Development

Format and test everything with:

```powershell
gofmt -w .
go test ./...
```

The project targets Go 1.25+.

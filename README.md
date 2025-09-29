# tscgit

`tscgit` is an interactive command-line coach for practicing Git workflows. It ships with ready-to-use lessons and a Bubble Tea-powered verification UI that guides students through checks in style.

## Quick Start

**Students: Choose your installation method**

### 🚀 One-Line Install (Recommended)

**Windows (PowerShell)**:
```powershell
iwr -useb https://raw.githubusercontent.com/rohit746/tscgit/main/install.ps1 | iex
```

**Linux/macOS**:
```bash
curl -fsSL https://raw.githubusercontent.com/rohit746/tscgit/main/install.sh | sh
```

### 📦 Other Installation Methods

**Go Users**:
```bash
go install github.com/rohit746/tscgit/cmd/tscgit@latest
```

**Download Pre-built Binaries**:
1. Visit the [releases page](https://github.com/rohit746/tscgit/releases)
2. Download the appropriate file for your system:
   - Windows: `tscgit_*_windows_amd64.zip`
   - macOS: `tscgit_*_darwin_amd64.tar.gz` (Intel) or `tscgit_*_darwin_arm64.tar.gz` (Apple Silicon)
   - Linux: `tscgit_*_linux_amd64.tar.gz`
3. Extract and place the binary in your `PATH`

**Package Managers**:
```bash
# Homebrew (macOS/Linux)
brew install rohit746/tap/tscgit

# APT (Ubuntu/Debian) - coming soon
# sudo apt install tscgit

# RPM (CentOS/RHEL/Fedora) - coming soon
# sudo dnf install tscgit
```

### ✅ Verify Installation

```bash
tscgit version
tscgit lessons
```

## System Requirements

- **Operating System**: Windows 10+, macOS 10.15+, or Linux (any modern distribution)
- **Architecture**: x86_64 (amd64) or ARM64
- **Dependencies**: None (statically compiled binary)
- **Git**: Required for lessons (install from [git-scm.com](https://git-scm.com))

## 🎓 Getting Started

After installation, start your Git learning journey:

### 1. See What's Available
```bash
tscgit lessons
```

### 2. Try Your First Lesson
```bash
# Start with environment setup
tscgit run 0

# Or jump into Git basics
tscgit verify init-basics
```

### 3. Practice in Your Own Repository
```bash
# Create a practice repository
mkdir git-practice && cd git-practice
git init

# Run verification lessons
tscgit verify init-basics
```

## Usage

**List available lessons**:
```bash
tscgit lessons
```

**Verify your Git skills** (interactive UI with real-time checks):
```bash
tscgit verify init-basics
tscgit verify -path /path/to/repo lesson-name
```

**Run guided practice scripts** (step-by-step command validation):
```bash
tscgit run 0    # Test your setup
tscgit run 1    # Install Git  
tscgit run 2    # Configure Git identity
```

**Check version**:
```bash
tscgit version
```

### 🎮 Interactive Features

- **Real-time verification**: See your progress as you work
- **Colorful feedback**: Clear visual indicators for success/failure  
- **Keyboard shortcuts**: Press `q` or `Ctrl+C` to exit anytime
- **Cross-platform**: Works on Windows PowerShell, macOS Terminal, and Linux shells

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
| `10`| Commit classics.csv          | Confirms the latest commit is `C: add classics.csv` while you're still on `add_classics`. |
| `11a`| Prepare merge history       | Inspects the pre-merge graph for commits `A:` through `D:`. |
| `11b`| Merge add_classics          | Confirms merge commit `E:` with multiple parents. |
| `12a`| Branch off D                | Ensures `update_dune` exists and points to commit `D: update contents.md`. |
| `12b`| Rebase update_dune          | Checks that commits `H:` and `I:` were added and the branch rebased onto `main`. |
| `13a`| Overwrite titles.md         | Simulates an accidental overwrite on `update_dune` and confirms commit `J:`. |
| `13b`| Soft reset to I              | Ensures the accidental commit is removed while keeping `titles.md` staged. |
| `13c`| Hard reset titles.md         | Teaches using a hard reset to drop the staged overwrite and restore the file. |

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

## 🛠 Development

### Prerequisites
- Go 1.25+
- Git

### Building from Source
```bash
# Clone the repository
git clone https://github.com/rohit746/tscgit.git
cd tscgit

# Build and install
go install ./cmd/tscgit

# Or build for development
go build -o bin/tscgit ./cmd/tscgit
```

### Testing
```bash
# Format code
gofmt -w .

# Run tests
go test ./...

# Test cross-compilation
GOOS=windows GOARCH=amd64 go build ./cmd/tscgit
GOOS=darwin GOARCH=arm64 go build ./cmd/tscgit
GOOS=linux GOARCH=amd64 go build ./cmd/tscgit
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Add your lessons or improvements
4. Test thoroughly
5. Submit a pull request

See the [architecture documentation](.github/copilot-instructions.md) for details on adding new lessons and scripts.

## 📚 For Educators

`tscgit` is designed to help educators teach Git effectively:

- **Standardized lessons**: Consistent learning experience across all students
- **Automated verification**: No manual grading of Git exercises  
- **Progressive difficulty**: Lessons build upon each other
- **Cross-platform**: Works in any environment your students use
- **Extensible**: Easy to add custom lessons for your curriculum

### Classroom Setup
```bash
# Students can install in seconds
curl -fsSL https://raw.githubusercontent.com/rohit746/tscgit/main/install.sh | sh

# Or provide direct download links from releases page
# https://github.com/rohit746/tscgit/releases
```

## 🤝 Contributing

We welcome contributions! Whether you're:
- **Students**: Report bugs, suggest improvements
- **Educators**: Add lessons for your curriculum  
- **Developers**: Improve the codebase

See our [contribution guidelines](CONTRIBUTING.md) and [architecture docs](.github/copilot-instructions.md).

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions  
- [Bubbles](https://github.com/charmbracelet/bubbles) - Common UI components

---

**Happy Git learning! 🎉**

If you find `tscgit` helpful, please ⭐ star the repository and share it with other Git learners!

````

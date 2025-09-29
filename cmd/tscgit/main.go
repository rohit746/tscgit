package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rohit746/tscgit/internal/gitutil"
	"github.com/rohit746/tscgit/internal/lessons"
	runlesson "github.com/rohit746/tscgit/internal/run"
	runui "github.com/rohit746/tscgit/internal/ui/run"
	verifyui "github.com/rohit746/tscgit/internal/ui/verify"
)

// Version information set at build time via ldflags
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage()
		return 1
	}

	switch args[0] {
	case "help", "--help", "-h":
		printUsage()
		return 0
	case "version", "--version", "-v":
		printVersion()
		return 0
	case "lessons", "list":
		printLessons()
		return 0
	case "verify":
		return handleVerify(args[1:])
	case "run":
		return handleRun(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printUsage()
		return 1
	}
}

func handleVerify(args []string) int {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	cwd := fs.String("path", "", "path to repository (defaults to current directory)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fs.PrintDefaults()
			return 0
		}
		return 1
	}
	remaining := fs.Args()
	if len(remaining) == 0 {
		fmt.Fprintln(os.Stderr, "verify requires a lesson ID. Try 'tscgit lessons' to list available options.")
		return 1
	}

	lessonID := remaining[0]
	lesson, err := lessons.Get(lessonID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	repo, err := gitutil.Open(context.Background(), *cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open git repository: %v\n", err)
		return 1
	}

	model := verifyui.NewModel(lesson, repo)
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "verification UI failed: %v\n", err)
		return 1
	}

	if model.AllPassed() {
		return 0
	}
	return 2
}

func handleRun(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "run requires a script ID. Try 'tscgit lessons' to list available options.")
		return 1
	}

	scriptID := args[0]
	script, ok := runlesson.Get(scriptID)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown run script: %s\n", scriptID)
		return 1
	}

	model := runui.NewModel(script)
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run UI failed: %v\n", err)
		return 1
	}

	if model.AllPassed() {
		return 0
	}
	return 2
}

func printUsage() {
	fmt.Fprintf(os.Stdout, `tscgit is a Git practice companion.

Usage:
  tscgit lessons             List available lessons
  tscgit verify <lesson-id>  Verify lesson progress in the current repo
  tscgit run <script-id>     Run terminal practice scripts
  tscgit version             Show version information

Flags:
  tscgit verify [-path DIR] <lesson-id>

`)
}

func printVersion() {
	fmt.Fprintf(os.Stdout, "tscgit version %s\n", Version)
	fmt.Fprintf(os.Stdout, "commit: %s\n", Commit)
	fmt.Fprintf(os.Stdout, "built: %s\n", Date)
}

func printLessons() {
	fmt.Fprintf(os.Stdout, "Verification lessons:\n\n")
	for _, lesson := range lessons.List() {
		fmt.Fprintf(os.Stdout, "  %-20s %s\n", lesson.ID, lesson.Title)
		if desc := strings.TrimSpace(lesson.Description); desc != "" {
			fmt.Fprintf(os.Stdout, "    %s\n", desc)
		}
	}

	fmt.Fprintf(os.Stdout, "\nRun lessons:\n\n")
	for _, script := range runlesson.List() {
		fmt.Fprintf(os.Stdout, "  %-20s %s\n", script.ID, script.Title)
		if desc := strings.TrimSpace(script.Description); desc != "" {
			fmt.Fprintf(os.Stdout, "    %s\n", desc)
		}
	}
}

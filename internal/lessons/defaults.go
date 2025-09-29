package lessons

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rohit746/tscgit/internal/gitutil"
)

func init() {
	Must(Register(lessonInitBasics()))
	Must(Register(lessonBranchBasics()))
}

func lessonInitBasics() *Lesson {
	return &Lesson{
		ID:          "init-basics",
		Title:       "Initialize a Git repository",
		Description: "Make your first commit, add a README, and practice writing meaningful messages.",
		Checks: []Check{
			{
				ID:          "first-commit",
				Title:       "Create at least one commit",
				Description: "Use git add and git commit so HEAD has history.",
				Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
					count, err := repo.CommitCount(ctx)
					if err != nil {
						return CheckResult{Err: err}
					}
					if count == 0 {
						return CheckResult{Passed: false, Message: "No commits detected. Run git commit to create your first snapshot."}
					}
					return CheckResult{Passed: true, Message: fmt.Sprintf("Great! You have %d commit(s) so far.", count)}
				},
			},
			{
				ID:          "readme-exists",
				Title:       "Add a README.md",
				Description: "Document what this practice repository is about.",
				Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
					exists, err := repo.FileExists("README.md")
					if err != nil {
						return CheckResult{Err: err}
					}
					if !exists {
						return CheckResult{Passed: false, Message: "Couldn't find README.md at the repo root."}
					}
					return CheckResult{Passed: true, Message: "README.md found—nice documentation!"}
				},
			},
			{
				ID:          "commit-message",
				Title:       "Write a descriptive commit message",
				Description: "Make sure your latest commit message is at least 5 characters long.",
				Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
					msg, err := repo.LastCommitMessage(ctx)
					if err != nil {
						return CheckResult{Err: err}
					}
					trimmed := strings.TrimSpace(msg)
					if len([]rune(trimmed)) < 5 {
						return CheckResult{Passed: false, Message: "Try writing a longer commit message that explains the change."}
					}
					return CheckResult{Passed: true, Message: fmt.Sprintf("Last commit message looks good: %q", trimmed)}
				},
			},
		},
	}
}

func lessonBranchBasics() *Lesson {
	return &Lesson{
		ID:          "branch-basics",
		Title:       "Practice branching",
		Description: "Create a feature branch, commit work on it, and keep main clean.",
		Checks: []Check{
			{
				ID:          "branch-exists",
				Title:       "Create feature/lesson-branch",
				Description: "Use git branch or git switch -c to create the feature branch.",
				Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
					exists, err := repo.HasBranch(ctx, "feature/lesson-branch")
					if err != nil {
						return CheckResult{Err: err}
					}
					if !exists {
						return CheckResult{Passed: false, Message: "Branch feature/lesson-branch not found."}
					}
					return CheckResult{Passed: true, Message: "feature/lesson-branch exists."}
				},
			},
			{
				ID:          "branch-current",
				Title:       "Check out the feature branch",
				Description: "Switch to feature/lesson-branch before committing work.",
				Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
					branch, err := repo.CurrentBranch(ctx)
					if err != nil {
						return CheckResult{Err: err}
					}
					if branch != "feature/lesson-branch" {
						return CheckResult{Passed: false, Message: fmt.Sprintf("Currently on %s. Switch to feature/lesson-branch.", branch)}
					}
					return CheckResult{Passed: true, Message: "Nice—you're working on the feature branch."}
				},
			},
			{
				ID:          "branch-commit",
				Title:       "Commit work on the feature branch",
				Description: "Create at least one commit that is ahead of main.",
				Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
					ahead, err := repo.CommitsAhead(ctx, "main", "feature/lesson-branch")
					if err != nil {
						return CheckResult{Err: err}
					}
					if ahead == 0 {
						return CheckResult{Passed: false, Message: "No commits found on feature/lesson-branch that aren't on main."}
					}
					return CheckResult{Passed: true, Message: fmt.Sprintf("feature/lesson-branch is ahead of main by %d commit(s).", ahead)}
				},
			},
			{
				ID:          "commit-message-tag",
				Title:       "Tag your commit message",
				Description: "Mention [branch] in the latest commit subject to flag branch work.",
				Verify: func(ctx context.Context, repo *gitutil.Repository) CheckResult {
					msg, err := repo.LastCommitMessage(ctx)
					if err != nil {
						return CheckResult{Err: err}
					}
					if !strings.Contains(msg, "[branch]") {
						return CheckResult{Passed: false, Message: "Add [branch] to your latest commit message for visibility."}
					}
					return CheckResult{Passed: true, Message: "[branch] tag detected in your latest commit message."}
				},
			},
		},
	}
}

// TimeoutContext wraps context.Background with a reasonable timeout for git calls used during verification.
func TimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

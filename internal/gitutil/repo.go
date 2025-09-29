package gitutil

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Repository encapsulates helpers for inspecting a local Git repository.
type Repository struct {
	Root string
}

// Open discovers the git repository that contains path. If path is empty, the
// current working directory is used.
func Open(ctx context.Context, path string) (*Repository, error) {
	if path == "" {
		p, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("gitutil: get working directory: %w", err)
		}
		path = p
	}

	out, err := gitCommand(ctx, path, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("gitutil: %w", err)
	}

	root := strings.TrimSpace(out)
	if root == "" {
		return nil, errors.New("gitutil: repository root is empty")
	}

	return &Repository{Root: root}, nil
}

// CurrentBranch returns the currently checked out branch.
func (r *Repository) CurrentBranch(ctx context.Context) (string, error) {
	out, err := r.git(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(out), err
}

// HasBranch reports whether a branch with the provided name exists.
func (r *Repository) HasBranch(ctx context.Context, name string) (bool, error) {
	out, err := r.git(ctx, "branch", "--list", name)
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		if strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(scanner.Text(), "*")), name) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

// FileExists reports whether the repository contains the given file path.
func (r *Repository) FileExists(relPath string) (bool, error) {
	_, err := os.Stat(filepath.Join(r.Root, relPath))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// CommitCount returns the number of commits reachable from HEAD.
func (r *Repository) CommitCount(ctx context.Context) (int, error) {
	out, err := r.git(ctx, "rev-list", "--count", "HEAD")
	if err != nil {
		if IsNoCommits(err) {
			return 0, nil
		}
		return 0, err
	}
	var count int
	_, scanErr := fmt.Sscanf(strings.TrimSpace(out), "%d", &count)
	if scanErr != nil {
		return 0, fmt.Errorf("gitutil: parse commit count: %w", scanErr)
	}
	return count, nil
}

// HasRemote returns true if the provided remote exists.
func (r *Repository) HasRemote(ctx context.Context, name string) (bool, error) {
	out, err := r.git(ctx, "remote")
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == name {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}

// LastCommitMessage returns the subject line of the most recent commit.
func (r *Repository) LastCommitMessage(ctx context.Context) (string, error) {
	out, err := r.git(ctx, "log", "-1", "--pretty=%s")
	if err != nil && IsNoCommits(err) {
		return "", nil
	}
	return strings.TrimSpace(out), err
}

// CommitsAhead reports how many commits compare has that are not in base.
func (r *Repository) CommitsAhead(ctx context.Context, base, compare string) (int, error) {
	if base == "" || compare == "" {
		return 0, errors.New("gitutil: base and compare branches are required")
	}
	spec := fmt.Sprintf("%s..%s", base, compare)
	out, err := r.git(ctx, "rev-list", "--count", spec)
	if err != nil {
		return 0, err
	}
	var count int
	if _, scanErr := fmt.Sscanf(strings.TrimSpace(out), "%d", &count); scanErr != nil {
		return 0, fmt.Errorf("gitutil: parse commits ahead count: %w", scanErr)
	}
	return count, nil
}

func (r *Repository) git(ctx context.Context, args ...string) (string, error) {
	return gitCommand(ctx, r.Root, args...)
}

func gitCommand(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git command failed: git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

// IsNoCommits reports whether the error corresponds to a repository without commits.
func IsNoCommits(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "does not have any commits yet") || strings.Contains(msg, "unknown revision or path not in the working tree")
}
